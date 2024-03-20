package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
)

func Seed(ctx context.Context, db *sql.DB, seedDir fs.FS) error {
	err := fs.WalkDir(seedDir, ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() || filepath.Ext(path) != "sql" {
			return nil
		}

		sql, err := fs.ReadFile(seedDir, path)
		if err != nil {
			return fmt.Errorf("read seed file=%s: %w", path, err)
		}

		if _, err := db.ExecContext(ctx, string(sql)); err != nil {
			return fmt.Errorf("execute seed file=%s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("sqldb seed: %w")
	}

	return nil
}
