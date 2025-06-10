package database

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	postgresCodeUniqueViolation = "23505"
)

var ErrUniqueViolation = errors.New("unique")

func IsUniqueViolation(err error) bool {
	return errors.Is(err, ErrUniqueViolation)
}

func TransformSQLError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == postgresCodeUniqueViolation {
			return ErrUniqueViolation
		}
	}

	return err
}
