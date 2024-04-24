package sqldb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  string
}

func (c Config) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Options,
	)
}

func ConnectPostgres(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("pgx connect to postgress: %w", err)
	}

	return pool, nil
}
