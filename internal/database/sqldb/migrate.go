package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, db *sql.DB, migrationDir fs.FS) error {
	provider, err := goose.NewProvider("postgres", db, migrationDir)
	if err != nil {
		return fmt.Errorf("create goose migration provider: %w", err)
	}

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
