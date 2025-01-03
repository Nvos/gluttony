package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"gluttony/tools/linter"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	const golangCiLintVersion = "1.61.0"
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("get work dir: %v", err))
	}

	binDir := filepath.Join(workDir, ".bin")
	if err := os.MkdirAll(binDir, 0644); err != nil {
		panic(fmt.Sprintf("create .bin directory: %v", err))
	}

	binFS := os.DirFS(binDir)

	tools := &Tools{
		logger:  logger,
		workDir: workDir,
		binFS:   binFS,
		binPath: binDir,
		linter:  linter.NewLinter(golangCiLintVersion, binDir, logger),
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "lint",
				Action: func(ctx *cli.Context) error {
					if err := tools.linter.Run(ctx.Context); err != nil {
						return err
					}

					return nil
				},
			},
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
						return err
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
	binFS   fs.FS
	binPath string

	linter *linter.Linter
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
