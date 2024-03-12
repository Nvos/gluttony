package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/config"
)

func ConnectPostgres(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, postgresURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("pgx connect to postgress: %w", err)
	}

	return pool, nil
}
