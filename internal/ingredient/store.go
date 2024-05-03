package ingredient

import (
	"context"
)

type Store interface {
	All(ctx context.Context, input AllIngredientsInput) ([]Ingredient, error)
	Single(ctx context.Context, input SingleInput) (Ingredient, error)
	Create(ctx context.Context, ingredient CreateIngredientInput) error
}
