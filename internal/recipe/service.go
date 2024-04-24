package recipe

import (
	"context"
	"fmt"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
	"gluttony/internal/x/assert"
)

type Service struct {
	beginner transaction.Beginner
	store    Store
}

func NewService(beginner transaction.Beginner, store Store) *Service {
	assert.Assert(beginner == nil, "transaction beginner is nil")
	assert.Assert(store == nil, "store is nil")

	return &Service{
		beginner: beginner,
		store:    store,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, recipe CreateRecipe) (int32, error) {
	tx, err := s.beginner.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}

	txStore, err := s.store.UnderTransaction(tx)
	if err != nil {
		return 0, fmt.Errorf("transactional store: %w", err)
	}

	fn := func() (int32, error) {
		recipeID, err := txStore.Create(ctx, recipe)
		if err != nil {
			return 0, fmt.Errorf("create recipe: %w", err)
		}

		insertableEdges := make([]IngredientEdge, 0, len(recipe.IngredientIDs))
		for i := range recipe.IngredientIDs {
			insertableEdges = append(insertableEdges, IngredientEdge{
				RecipeID:     recipeID,
				IngredientID: recipe.IngredientIDs[i],
			})
		}

		if err := txStore.CreateIngredientEdges(ctx, insertableEdges); err != nil {
			return 0, fmt.Errorf("create ingredient edges: %w", err)
		}

		return recipeID, nil
	}

	recipeID, err := fn()
	if err := transaction.ResolveTx(ctx, err, tx); err != nil {
		return 0, fmt.Errorf("resolve create recipe tx: %w", err)
	}

	return recipeID, nil
}

func (s *Service) SingleRecipe(ctx context.Context, locale i18n.Locale, recipeID int32) (FullRecipe, error) {
	single, err := s.store.Single(ctx, locale, recipeID)
	if err != nil {
		return FullRecipe{}, fmt.Errorf("single recipe: %w", err)
	}

	return single, nil
}

func (s *Service) AllRecipes(ctx context.Context, input AllRecipesInput) ([]Recipe, error) {
	all, err := s.store.All(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("all paged recipes: %w", err)
	}

	return all, nil
}
