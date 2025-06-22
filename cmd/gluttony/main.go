package main

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

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("get work dir path: %v", err))
	}

	cfg, err := config.NewConfig(wd)
	if err != nil {
		fmt.Printf("Create config failed: '%v', Aborting startup\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	a, err := app.New(cfg)
	if err != nil {
		fmt.Printf("Create new app failed: '%v', Aborting startup\n", err)
		os.Exit(1)
	}

	if err := a.Run(groupCtx, group); err != nil {
		fmt.Printf("Run new app failed: '%v', Aborting startup\n", err)
		os.Exit(1)
	}

	if err := group.Wait(); err != nil {
		fmt.Printf("Run new app failed: '%v', Aborting startup\n", err)
		os.Exit(1)
	}
}
