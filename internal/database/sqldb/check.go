package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsDBRunning(ctx context.Context, db *sql.DB) (bool, error) {
	err := db.PingContext(ctx)
	var connErr *pgconn.ConnectError
	if errors.As(err, &connErr) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("IsDBRunning: %w", err)
	}

	return true, nil
}
