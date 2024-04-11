package sqldb

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("not found")

type UniqueConstraintError struct {
	Column string
}

func (u UniqueConstraintError) Error() string {
	return fmt.Sprintf("unique constraint failed for column %s", u.Column)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func AsUniqueConstraint(err error) (*UniqueConstraintError, bool) {
	var target *UniqueConstraintError
	if errors.As(err, &target) {
		return target, true
	}

	return nil, false
}

func TransformPgxError(err error) error {
	var target *pgconn.PgError
	if errors.As(err, &target) {
		switch target.Code {
		case "23505":
			return &UniqueConstraintError{
				Column: target.ColumnName,
			}
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}
