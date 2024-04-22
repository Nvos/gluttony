package ingredient

import (
	"context"
)

type Store interface {
	All(ctx context.Context, input AllIngredientsInput) ([]Ingredient, error)
	Create(ctx context.Context, ingredient CreateIngredientInput) error
}
