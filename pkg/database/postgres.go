package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
	"net/url"
)

type Config struct {
	Name     string
	User     string
	Host     string
	Port     int
	Password string
}

func (c *Config) ConnectionURL() string {
	u := &url.URL{
		Scheme: "postgres",
		Host:   c.Host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	q.Add("sslmode", "disable")

	u.RawQuery = q.Encode()

	return u.String()
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("new db: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

func Migrate(pool *pgxpool.Pool, files fs.FS) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(files)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	return nil
}
