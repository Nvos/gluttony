package admin

import (
	"context"
	"fmt"
	config2 "gluttony/config"
	"gluttony/migrations"
	"gluttony/x/sqlx"
)

func RunMigrations(ctx context.Context, cfg *config2.Config, sec *config2.Secret) error {
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
