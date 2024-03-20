package sqldb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/config"
)

func postgresURL(cfg config.Database) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=public&sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func ConnectPostgres(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, postgresURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("pgx connect to postgress: %w", err)
	}

	return pool, nil
}
