package recipe

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/database/pagination"
	"gluttony/internal/database/transaction"
	"gluttony/internal/recipe/postgresql"
)

var _ Store = (*PostgresStore)(nil)

type PostgresStore struct {
	pool    postgresql.DBTX
	queries *postgresql.Queries
}

func (s *PostgresStore) UnderTransaction(tx transaction.Transaction) (Store, error) {
	pgxTx, err := transaction.GetPgxTx(tx)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		pool:    pgxTx,
		queries: s.queries.WithTx(pgxTx),
	}, nil
}

func (s *PostgresStore) Create(ctx context.Context, value CreateRecipe) (int32, error) {
	recipeId, err := s.queries.CreateRecipe(ctx, postgresql.CreateRecipeParams{
		Description: value.Description,
		Name:        value.Name,
	})
	if err != nil {
		return 0, fmt.Errorf("postgres: create recipe: %w", err)
	}

	return recipeId, nil
}

func (s *PostgresStore) CreateRecipeSteps(ctx context.Context, recipeID int32, steps []CreateStep) error {
	if recipeID <= 0 {
		return fmt.Errorf("invalid recipe id=%d", recipeID)
	}

	if len(steps) == 0 {
		return errors.New("create recipe steps list is empty")
	}

	dbSteps := make([]postgresql.CreateRecipeStepsParams, 0, len(steps))
	for i := range steps {
		dbSteps = append(dbSteps, postgresql.CreateRecipeStepsParams{
			RecipeID:    recipeID,
			Description: steps[i].Description,
			Order:       steps[i].Order,
		})
	}

	createdCount, err := s.queries.CreateRecipeSteps(ctx, dbSteps)
	if err != nil {
		return fmt.Errorf("postgres: create recipe steps: %w", err)
	}

	if len(steps) != int(createdCount) {
		return fmt.Errorf("postgres: not all recipe steps were persisted")
	}

	return nil
}

func (s *PostgresStore) Single(ctx context.Context, id int32) (FullRecipe, error) {
	rows, err := s.queries.SingleRecipe(ctx, id)
	if err != nil {
		return FullRecipe{}, fmt.Errorf("postgres: single recipe by id=%d: %w", id, err)
	}

	if len(rows) == 0 {
		return FullRecipe{}, fmt.Errorf("postgres: single recipe by id=%d not found", id)
	}

	steps := make([]Step, 0, len(rows))
	for i := range rows {
		step := Step{
			ID:          rows[i].RecipeStep.ID,
			Order:       rows[i].RecipeStep.Order,
			Description: rows[i].RecipeStep.Description,
		}

		steps = append(steps, step)
	}

	recipe := FullRecipe{
		Recipe: Recipe{
			ID:          rows[0].Recipe.ID,
			Name:        rows[0].Recipe.Name,
			Description: rows[0].Recipe.Description,
		},
		Steps: steps,
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
