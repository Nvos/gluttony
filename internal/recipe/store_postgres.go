package recipe

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/database/pagination"
	"gluttony/internal/recipe/postgresql"
)

var _ Store = (*PostgresStore)(nil)

type PostgresStore struct {
	pool    *pgxpool.Pool
	queries *postgresql.Queries
}

func (s *PostgresStore) withTx(
	ctx context.Context,
	cb func(ctx context.Context, tx pgx.Tx, queries *postgresql.Queries) (int32, error),
) (id int32, err error) {

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin postgres transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); err != nil {
			err = errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
		}
	}()

	id, err = cb(ctx, tx, s.queries.WithTx(tx))
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return id, nil
}

func (s *PostgresStore) Create(ctx context.Context, value CreateRecipe) (int32, error) {
	recipeId, err := s.withTx(ctx, func(ctx context.Context, tx pgx.Tx, queries *postgresql.Queries) (int32, error) {
		recipeId, err := queries.CreateRecipe(ctx, postgresql.CreateRecipeParams{
			Description: value.Description,
			Name:        value.Name,
		})
		if err != nil {
			return 0, fmt.Errorf("postgres: create recipe: %w", err)
		}

		if len(value.Steps) == 0 {
			return recipeId, nil
		}

		steps := make([]postgresql.CreateRecipeStepsParams, 0, len(value.Steps))
		for i := range value.Steps {
			steps = append(steps, postgresql.CreateRecipeStepsParams{
				RecipeID:    recipeId,
				Description: value.Steps[i].Description,
				Order:       value.Steps[i].Order,
			})
		}

		createdCount, err := queries.CreateRecipeSteps(ctx, steps)
		if err != nil {
			return 0, fmt.Errorf("postgres: create recipe steps: %w", err)
		}

		if len(steps) != int(createdCount) {
			return 0, fmt.Errorf("postgres: not all recipe steps were persisted")
		}

		return recipeId, nil
	})

	if err != nil {
		return 0, fmt.Errorf("postgres: create recipe: %w", err)
	}

	return recipeId, nil
}

func (s *PostgresStore) Single(ctx context.Context, id int32) (Recipe, error) {
	rows, err := s.queries.SingleRecipe(ctx, id)
	if err != nil {
		return Recipe{}, fmt.Errorf("postgres: single recipe by id=%d: %w", id, err)
	}

	if len(rows) == 0 {
		return Recipe{}, fmt.Errorf("postgres: single recipe by id=%d not found", id)
	}

	recipe := Recipe{
		ID:          rows[0].Recipe.ID,
		Name:        rows[0].Recipe.Name,
		Description: rows[0].Recipe.Description,
	}

	steps := make([]Step, 0, len(rows))
	for i := range rows {
		step := Step{
			ID:          steps[i].ID,
			Order:       steps[i].Order,
			Description: steps[i].Description,
		}

		steps = append(steps, step)
	}

	return recipe, nil
}

func (s *PostgresStore) All(ctx context.Context, search string, pagination pagination.OffsetPagination) ([]Recipe, error) {
	params := postgresql.AllRecipesParams{
		Offset:             pagination.Offset,
		Limit:              pagination.Limit,
		WebsearchToTsquery: search,
	}

	rows, err := s.queries.AllRecipes(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("postgres: all recipes: %w", err)
	}

	recipes := make([]Recipe, 0, len(rows))
	for i := range rows {
		recipe := Recipe{
			ID:          rows[i].ID,
			Name:        rows[i].Name,
			Description: rows[i].Description,
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func NewPostgresStore(pool *pgxpool.Pool) (*PostgresStore, error) {
	if pool == nil {
		return nil, fmt.Errorf("new postgres store: pgxpool is nil")
	}

	return &PostgresStore{
		pool:    pool,
		queries: postgresql.New(pool),
	}, nil
}
