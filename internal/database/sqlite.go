package database

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func New(workDir string) (*sql.DB, error) {
	dbPath := filepath.Join(workDir, "sqlite.db")
	create, err := os.Create(dbPath)
	if err != nil {
		return nil, err
	}

	if err := create.Close(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate sqlite db: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
