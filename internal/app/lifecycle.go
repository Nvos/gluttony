package app

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/user"
	"gluttony/pkg/livereload"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

func (app *App) Start(ctx context.Context, group *errgroup.Group) error {
	app.logger.InfoContext(ctx, "Directories", slog.String("root", app.cfg.WorkDirectoryPath))
	// TODO: move to cmd, left for now as it is convenient
	if err := app.userService.Create(ctx, user.CreateInput{
		Username: "admin",
		Password: "admin",
		Role:     user.RoleAdmin,
	}); err != nil {
		return fmt.Errorf("create initial admin user: %w", err)
	}

	group.Go(func() error {
		if app.cfg.Mode == Prod {
			return nil
		}

		if err := app.liveReload.Watch(ctx, livereload.WatchConfig{
			Extensions: []string{".gohtml", ".html", ".css", ".js"},
			Directories: []string{
				filepath.Clean("assets"),
				filepath.Clean("web/templates"),
			},
		}); err != nil {
			return fmt.Errorf("start livereload watch: %w", err)
		}

		app.logger.InfoContext(ctx, "Shutting down live reload")
		return nil
	})

	group.Go(func() error {
		app.logger.InfoContext(ctx, "Starting HTTP server", slog.String("address", app.httpServer.Addr))

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
