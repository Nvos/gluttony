package livereload

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

type WatchConfig struct {
	Directories []string
	Extensions  []string
}

type LiveReload struct {
	cond   sync.Cond
	logger *slog.Logger
}

func New(logger *slog.Logger) *LiveReload {
	return &LiveReload{
		cond:   sync.Cond{L: &sync.Mutex{}},
		logger: logger,
	}
}

func (reload *LiveReload) Handle(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Block here until next reload event
	reload.cond.L.Lock()
	reload.cond.Wait()
	reload.cond.L.Unlock()

	_, _ = fmt.Fprintf(w, "data: reload\n\n")
	w.(http.Flusher).Flush()
}

func (reload *LiveReload) Watch(ctx context.Context, cfg WatchConfig) error {
	if cfg.Directories == nil {
		return fmt.Errorf("no watch directories provided")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create new directory watcher: %w", err)
	}

	for i := range cfg.Directories {
		directories, err := collectDirectories(cfg.Directories[i])
		if err != nil {
			return err
		}

		for _, directory := range directories {
			if err := w.Add(directory); err != nil {
				return fmt.Errorf("add watch directory: %w", err)
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
				return fmt.Errorf("watch handle create: %w", err)
			}

			debounce(reload.cond.Broadcast)
		case e.Has(fsnotify.Write):
			debounce(reload.cond.Broadcast)
		case e.Has(fsnotify.Remove), e.Has(fsnotify.Rename):
			directories, err := collectDirectories(e.Name)
			if err != nil {
				return fmt.Errorf("watch handle remove/rename: %w", err)
			}

			for _, v := range directories {
				if err := w.Remove(v); err != nil {
					reload.logger.WarnContext(
						ctx,
						"Remove watch directory",
						slog.String("err", err.Error()),
					)
				}
			}

			if err := w.Remove(e.Name); err != nil {
				reload.logger.WarnContext(
					ctx,
					"Remove watch file",
					slog.String("err", err.Error()),
				)
			}
		}

		return nil
	}

	defer w.Close()

	for {
		select {
		case <-ctx.Done():
			// Free SSE lock
			reload.cond.Broadcast()
			return nil
		case err := <-w.Errors:
			if err != nil {
				reload.logger.ErrorContext(
					ctx,
					"File watcher error",
					slog.String("err", err.Error()),
				)
			}
		case e := <-w.Events:
			if err := handleEvent(e); err != nil {
				reload.logger.WarnContext(
					ctx,
					"File watcher handle event failed",
					slog.String("err", err.Error()),
				)
			}
		}
	}
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
