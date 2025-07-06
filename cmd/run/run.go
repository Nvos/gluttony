package run

import (
	"context"
	"fmt"
	"gluttony/internal/app"
	"gluttony/internal/config"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func Run(rootCtx context.Context, cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(rootCtx,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	a, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("create new app: %w", err)
	}

	if err := a.Run(groupCtx, group); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("wait for app: %w", err)
	}

	return nil
}
