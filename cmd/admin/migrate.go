package admin

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/migrations"
	"gluttony/x/sqlx"
)

func RunMigrations(ctx context.Context, cfg *config.Config, sec *config.Secret) error {
	pool, err := sqlx.New(ctx, cfg.Database, sec.Database)
	if err != nil {
		return fmt.Errorf("create db: %w", err)
	}
	defer pool.Close()

	if err := sqlx.Migrate(pool, migrations.Migrations); err != nil {
		return fmt.Errorf("migrate db: %w", err)
	}

	return nil
}
