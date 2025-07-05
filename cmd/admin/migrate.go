package admin

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/migrations"
	"gluttony/pkg/database"
)

func RunMigrations(ctx context.Context) error {
	cfg, err := config.NewConfig()
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
