package recipe

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/database/pagination"
	"gluttony/internal/database/transaction"
)

type Service struct {
	beginner transaction.Beginner
	store    Store
}

func NewService(beginner transaction.Beginner, store Store) (*Service, error) {
	if beginner == nil {
		return nil, errors.New("recipe service beginner is nil")
	}

	if store == nil {
		return nil, errors.New("recipe service store is nil")
	}

	return &Service{
		beginner: beginner,
		store:    store,
	}, nil
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

		if err := txStore.CreateRecipeSteps(ctx, recipeID, recipe.Steps); err != nil {
			return 0, fmt.Errorf("create recipe steps: %w", err)
		}

		return recipeID, nil
	}

	recipeID, err := fn()
	if err := transaction.ResolveTx(ctx, err, tx); err != nil {
		return 0, fmt.Errorf("resolve create recipe tx: %w", err)
	}

	return recipeID, nil
}

func (s *Service) SingleRecipe(ctx context.Context, recipeID int32) (FullRecipe, error) {
	single, err := s.store.Single(ctx, recipeID)
	if err != nil {
		return FullRecipe{}, fmt.Errorf("single recipe: %w", err)
	}

	return single, nil
}

func (s *Service) AllRecipes(ctx context.Context, search string, params pagination.OffsetPagination) ([]Recipe, error) {
	if err := pagination.ValidateOffsetPagination(params); err != nil {
		return nil, err
	}

	all, err := s.store.All(ctx, search, params)
	if err != nil {
		return nil, fmt.Errorf("all paged recipes: %w", err)
	}

	return all, nil
}
