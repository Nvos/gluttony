package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/exec"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("get work dir: %v", err))
	}

	tools := &Tools{
		logger:  logger,
		workDir: workDir,
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "tailwind",
				Action: func(ctx *cli.Context) error {
					group := &errgroup.Group{}

					group.Go(func() error {
						if err := tools.liveTailwind(ctx.Context); err != nil {
							return fmt.Errorf("live tailwind: %w", err)
						}

						return nil
					})

					if err := group.Wait(); err != nil {
						return fmt.Errorf("tailwind wait: %w", err)
					}

					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(fmt.Sprintf("cmd failed: %v", err))
	}
}

type Tools struct {
	logger *slog.Logger

	workDir string
}

func (t *Tools) liveTailwind(ctx context.Context) error {
	return t.runCmd(
		ctx,
		"pnpm",
		"tailwind:watch",
	)
}

func (t *Tools) runCmd(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(
		ctx,
		name,
		args...,
	)

	cmd.Dir = t.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	t.logger.InfoContext(ctx, "Run command", slog.String("cmd", cmd.String()))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run cmd: %w", err)
	}

	return nil
}
