package app

import (
	"context"
	"errors"
	"fmt"
	"gluttony/config"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"time"
)

const defaultHttpServerTimeout = time.Second * 15

func (app *App) Run(ctx context.Context, group *errgroup.Group) error {
	app.logger.InfoContext(
		ctx,
		"Configuration",
		slog.String("mode", string(app.cfg.Mode)),
		slog.String("rootDir", app.cfg.WorkDir),
	)

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
	shutdownTimeout := defaultHttpServerTimeout
	if app.cfg.Mode == config.ModeDev {
		shutdownTimeout = 0
	}
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
