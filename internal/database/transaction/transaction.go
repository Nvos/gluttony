package transaction

import (
	"context"
	"errors"
	"fmt"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Beginner interface {
	Begin(ctx context.Context) (Transaction, error)
}

func ResolveTx(ctx context.Context, err error, tx Transaction) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("rollback tx: %w: %w", rollbackErr, err))
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
