package transaction

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPgxTx(tx Transaction) (pgx.Tx, error) {
	out, ok := tx.(pgx.Tx)
	if !ok {
		return nil, errors.New("expected pgx transaction")
	}

	return out, nil
}

type PgxBeginner struct {
	pool *pgxpool.Pool
}

func (b *PgxBeginner) Begin(ctx context.Context) (Transaction, error) {
	return b.pool.Begin(ctx)
}

func NewPgxBeginner(pool *pgxpool.Pool) Beginner {
	return &PgxBeginner{
		pool: pool,
	}
}
