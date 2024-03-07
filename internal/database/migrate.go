package database

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(cfg Config) (err error) {
	dir, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("resolve migration directory: %w", err)
	}

	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			dir,
		),
	)
	if err != nil {
		return fmt.Errorf("load migration directory: %w", err)
	}

	defer func() {
		if closeErr := workdir.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close migration workdir: %w", closeErr))
		}
	}()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return fmt.Errorf("initialize atlas migration client: %w", err)
	}

	_, err = client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: postgresURL(cfg),
	})
	if err != nil {
		return fmt.Errorf("apply atlas migrations: %w", err)
	}

	return nil
}
