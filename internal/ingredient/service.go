package ingredient

import (
	"context"
	"fmt"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Single(ctx context.Context, input SingleInput) (Ingredient, error) {
	single, err := s.store.Single(ctx, input)
	if err != nil {
		return Ingredient{}, fmt.Errorf("single ingredient by id=%d: %w", input.ID, err)
	}

	return single, nil
}

func (s *Service) All(ctx context.Context, input AllIngredientsInput) ([]Ingredient, error) {
	all, err := s.store.All(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("all ingredients: %w", err)
	}

	return all, nil
}

func (s *Service) Create(ctx context.Context, input CreateIngredientInput) error {
	if err := s.store.Create(ctx, input); err != nil {
		return fmt.Errorf("create ingredient: %w", err)
	}

	return nil
}
