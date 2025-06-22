package commands

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/migrations"
	"gluttony/pkg/database"
	"os"
)

func RunMigrations(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get work dir path: %w", err)
	}

	cfg, err := config.NewConfig(wd)
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	pool, err := database.New(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("create db: %w", err)
	}
	defer pool.Close()

	if err := database.Migrate(pool, migrations.Migrations); err != nil {
		return fmt.Errorf("migrate db: %w", err)
	}

	return nil
}
