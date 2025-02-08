package app

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/livereload"
	"golang.org/x/sync/errgroup"
	"net/http"
	"path/filepath"
	"time"
)

func (app *App) Start(ctx context.Context, group *errgroup.Group) error {
	// TODO: move to cmd, left for now as it is convenient
	if err := app.userService.Create(ctx, "admin", "admin"); err != nil {
		return err
	}

	group.Go(func() error {
		if err := app.liveReload.Watch(ctx, livereload.WatchConfig{
			Extensions: []string{".gohtml", ".html", ".css", ".js"},
			Directories: []string{
				filepath.Join("assets"),
				filepath.Clean(filepath.Join("internal/templating/templates")),
				filepath.Clean(filepath.Join("internal/user/templates")),
				filepath.Clean(filepath.Join("internal/recipe/templates")),
			},
		}); err != nil {
			return fmt.Errorf("start livereload watch: %w", err)
		}

		app.logger.InfoContext(ctx, "Shutting down live reload")
		return nil
	})

	group.Go(func() error {
		if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("start http server: %w", err)
		}

		app.logger.InfoContext(ctx, "Shutting down live http server")
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()

		app.logger.InfoContext(ctx, "Starting graceful shutdown")
		if err := app.stop(); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}

		app.logger.InfoContext(ctx, "Graceful shutdown completed")
		return nil
	})

	return nil
}

func (app *App) stop() error {
	const shutdownTimeout = 15 * time.Second
	shutdownCtx, cancelFn := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelFn()

	if err := app.recipeService.Stop(); err != nil {
		return fmt.Errorf("gracefully stopping recipe service: %w", err)
	}

	if err := app.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gracefully stopping http server: %w", err)
	}

	return nil
}
