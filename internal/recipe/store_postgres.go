package recipe

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
	"gluttony/internal/recipe/postgresql"
	"gluttony/internal/x/assert"
)

var _ Store = (*StorePostgres)(nil)

type StorePostgres struct {
	pool    postgresql.DBTX
	queries *postgresql.Queries
}

func (s *StorePostgres) CreateIngredientEdges(ctx context.Context, value []IngredientEdge) error {
	params := make([]postgresql.CreateRecipeIngredientEdgesParams, 0, len(value))
	for i := range value {
		params = append(params, postgresql.CreateRecipeIngredientEdgesParams{
			RecipeID:     value[i].RecipeID,
			IngredientID: value[i].IngredientID,
		})
	}

	if _, err := s.queries.CreateRecipeIngredientEdges(ctx, params); err != nil {
		return fmt.Errorf("create recipe ingredient edges: %w", err)
	}

	return nil
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

func (s *StorePostgres) Create(ctx context.Context, value CreateRecipe) (int32, error) {
	recipeId, err := s.queries.CreateRecipe(ctx, postgresql.CreateRecipeParams{
		Description: i18n.NewField(value.Locale, value.Description).JSONBytes(),
		Name:        i18n.NewField(value.Locale, value.Name).JSONBytes(),
		Content:     i18n.NewField(value.Locale, value.Content).JSONBytes(),
	})
	if err != nil {
		return 0, fmt.Errorf("postgres: create recipe: %w", err)
	}

	return recipeId, nil
}

func (s *StorePostgres) Single(ctx context.Context, locale i18n.Locale, id int32) (FullRecipe, error) {
	row, err := s.queries.SingleRecipe(ctx, postgresql.SingleRecipeParams{
		Locale:   string(locale),
		RecipeID: id,
	})
	if err != nil {
		return FullRecipe{}, fmt.Errorf("postgres: single recipe by id=%d: %w", id, err)
	}

	//if len(rows) == 0 {
	//	return FullRecipe{}, fmt.Errorf("postgres: single recipe by id=%d not found", id)
	//}
	//
	//steps := make([]Step, 0, len(rows))
	//for i := range rows {
	//	step := Step{
	//		ID:          rows[i].RecipeStep.ID,
	//		Order:       rows[i].RecipeStep.Order,
	//		Description: rows[i].RecipeStep.Description,
	//	}
	//
	//	steps = append(steps, step)
	//}

	recipe := FullRecipe{
		Recipe: Recipe{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
		},
	}

	return recipe, nil
}

func (s *StorePostgres) All(ctx context.Context, input AllRecipesInput) ([]Recipe, error) {
	params := postgresql.AllRecipesParams{
		Offset: input.Pagination.Offset,
		Limit:  input.Pagination.Limit,
		Search: input.Search,
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

func NewStorePostgres(pool *pgxpool.Pool) *StorePostgres {
	assert.Assert(pool != nil, "pgx pool is nil")

	return &StorePostgres{
		pool:    pool,
		queries: postgresql.New(pool),
	}
}
