package ingredient

import (
	"context"
	"fmt"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
	"gluttony/internal/ingredient/postgresql"
)

var _ Store = (*StorePostgres)(nil)

type StorePostgres struct {
	pool    postgresql.DBTX
	queries *postgresql.Queries
}

func NewStorePostgres(pool postgresql.DBTX) *StorePostgres {
	return &StorePostgres{
		pool:    pool,
		queries: postgresql.New(pool),
	}
}

func (s *StorePostgres) Create(ctx context.Context, input CreateIngredientInput) error {
	if err := s.queries.CreateIngredient(ctx, i18n.NewField(input.Locale, input.Name).JSONBytes()); err != nil {
		return fmt.Errorf("create ingredient: %w", err)
	}

	return nil
}

func (s *StorePostgres) All(
	ctx context.Context,
	input AllIngredientsInput,
) ([]Ingredient, error) {
	rows, err := s.queries.AllIngredients(ctx, postgresql.AllIngredientsParams{
		Locale: string(input.Locale),
		Offset: input.Pagination.Offset,
		Limit:  input.Pagination.Limit,
		Search: input.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres all ingredients: %w", err)
	}

	out := make([]Ingredient, 0, len(rows))
	for i := range rows {
		out = append(out, Ingredient{
			ID:   rows[i].ID,
			Name: rows[i].Name,
		})
	}

	return out, nil
}

func (s *StorePostgres) UnderTransaction(tx transaction.Transaction) (Store, error) {
	pgxTx, err := transaction.GetPgxTx(tx)
	if err != nil {
		return nil, err
	}

	return &StorePostgres{
		pool:    pgxTx,
		queries: s.queries.WithTx(pgxTx),
	}, nil
}