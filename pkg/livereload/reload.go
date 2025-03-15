package livereload

import (
	"context"
	"errors"
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
	cond   *sync.Cond
	logger *slog.Logger
}

type watcher struct {
	fsnotify   *fsnotify.Watcher
	debounceFn func(f func())
	cfg        WatchConfig
	logger     *slog.Logger
	cond       *sync.Cond
}

func (w *watcher) close() error {
	if err := w.fsnotify.Close(); err != nil {
		return fmt.Errorf("close fsnotify: %w", err)
	}

	return nil
}

func newWatcher(cfg WatchConfig, logger *slog.Logger, cond *sync.Cond) (*watcher, error) {
	const debounceTime = time.Millisecond * 100
	debounceFn := NewDebounce(debounceTime)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create new directory fsnotify: %w", err)
	}

	for i := range cfg.Directories {
		directories, err := collectDirectories(cfg.Directories[i])
		if err != nil {
			return nil, err
		}

		for _, directory := range directories {
			if err := w.Add(directory); err != nil {
				return nil, fmt.Errorf("add watch directory: %w", err)
			}
		}
	}

	return &watcher{
		debounceFn: debounceFn,
		cfg:        cfg,
		fsnotify:   w,
		logger:     logger,
		cond:       cond,
	}, nil
}

func New(logger *slog.Logger) *LiveReload {
	return &LiveReload{
		cond:   sync.NewCond(&sync.Mutex{}),
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
		return errors.New("no watch directories provided")
	}

	w, err := newWatcher(cfg, reload.logger, reload.cond)
	if err != nil {
		return fmt.Errorf("new livereload fsnotify: %w", err)
	}

	defer func() {
		_ = w.close()
	}()

	for {
		select {
		case <-ctx.Done():
			// Free SSE lock
			reload.cond.Broadcast()
			return nil
		case err := <-w.fsnotify.Errors:
			if err != nil {
				reload.logger.ErrorContext(
					ctx,
					"File fsnotify error",
					slog.String("err", err.Error()),
				)
			}
		case e := <-w.fsnotify.Events:
			if err := w.handleEvent(ctx, e); err != nil {
				reload.logger.WarnContext(
					ctx,
					"File fsnotify handle event failed",
					slog.String("err", err.Error()),
				)
			}
		}
	}
}

func (w *watcher) handleEvent(ctx context.Context, e fsnotify.Event) error {
	matchingExtension := false
	for i := range w.cfg.Extensions {
		if filepath.Ext(e.Name) == w.cfg.Extensions[i] {
			matchingExtension = true
			break
		}
	}

	if !matchingExtension {
		return nil
	}

	switch {
	case e.Has(fsnotify.Create):
		if err := w.fsnotify.Add(e.Name); err != nil {
			return fmt.Errorf("watch handle create: %w", err)
		}

		w.debounceFn(w.cond.Broadcast)
	case e.Has(fsnotify.Write):
		w.debounceFn(w.cond.Broadcast)
	case e.Has(fsnotify.Remove), e.Has(fsnotify.Rename):
		directories, err := collectDirectories(e.Name)
		if err != nil {
			return fmt.Errorf("watch handle remove/rename: %w", err)
		}

		for _, v := range directories {
			if err := w.fsnotify.Remove(v); err != nil {
				w.logger.WarnContext(
					ctx,
					"Remove watch directory",
					slog.String("err", err.Error()),
				)
			}
		}

		if err := w.fsnotify.Remove(e.Name); err != nil {
			w.logger.WarnContext(
				ctx,
				"Remove watch file",
				slog.String("err", err.Error()),
			)
		}
	}

	return nil
}

func collectDirectories(dirPath string) ([]string, error) {
	var directories []string

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return err
		}

		directories = append(directories, path)

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	return directories, nil
}
