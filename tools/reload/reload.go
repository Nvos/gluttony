package reload

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"io/fs"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

type WatchConfig struct {
	Directories []string
	Extensions  []string
}

type Reload struct {
	upgrader websocket.Upgrader
	cond     sync.Cond
}

func New() *Reload {
	return &Reload{
		cond:     sync.Cond{L: &sync.Mutex{}},
		upgrader: websocket.Upgrader{},
	}
}

func (reload *Reload) Handle(w http.ResponseWriter, r *http.Request) {

	conn, err := reload.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// todo: handle err
		panic(err)
		return
	}

	// Block here until next reload event
	reload.cond.L.Lock()
	reload.cond.Wait()
	reload.cond.L.Unlock()

	conn.WriteMessage(websocket.TextMessage, []byte("reload"))
	conn.Close()
}

func (reload *Reload) Watch(ctx context.Context, cfg WatchConfig) error {
	if cfg.Directories == nil {
		panic("reload: watch Directories not set")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("reload: create new watcher: %w", err)
	}

	for i := range cfg.Directories {
		directories, err := collectDirectories(cfg.Directories[i])
		if err != nil {
			return err
		}

		for _, directory := range directories {
			if err := w.Add(directory); err != nil {
				return fmt.Errorf("reload: add watch directory: %w", err)
			}
		}
	}

	debounce := NewDebounce(time.Millisecond * 100)

	handleEvent := func(e fsnotify.Event) error {
		matchingExtension := false
		for i := range cfg.Extensions {
			if filepath.Ext(e.Name) == cfg.Extensions[i] {
				matchingExtension = true
				break
			}
		}

		if !matchingExtension {
			return nil
		}

		switch {
		case e.Has(fsnotify.Create):
			if err := w.Add(e.Name); err != nil {
				return fmt.Errorf("reload: add watch directory: %w", err)
			}

			debounce(reload.cond.Broadcast)
		case e.Has(fsnotify.Write):
			debounce(reload.cond.Broadcast)
		case e.Has(fsnotify.Remove), e.Has(fsnotify.Rename):
			directories, err := collectDirectories(e.Name)
			if err != nil {
				return fmt.Errorf("reload: collect Directories: %w", err)
			}

			for _, v := range directories {
				w.Remove(v)
			}

			w.Remove(e.Name)
		}

		return nil
	}

	go func() {
		defer w.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-w.Errors:
				// todo: proper log
				if err != nil {
					println(err.Error())
				}
			case e := <-w.Events:
				if err := handleEvent(e); err != nil {
					// todo: proper log
					println(err.Error())
				}
			}
		}
	}()

	return nil
}

func collectDirectories(dirPath string) ([]string, error) {
	var directories []string

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil

		}

		directories = append(directories, path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return directories, nil
}
