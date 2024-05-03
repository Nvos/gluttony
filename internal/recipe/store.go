package recipe

import (
	"context"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
)

// TODO, 01/05/2024: fix name, move to recipe.go, add constructor
type IngredientEdge struct {
	RecipeID     int32
	IngredientID int32
	Amount       int32
	Count        int32
	Note         string
}

type Store interface {
	UnderTransaction(tx transaction.Transaction) (Store, error)
	Single(ctx context.Context, locale i18n.Locale, id int32) (FullRecipe, error)
	All(ctx context.Context, input AllRecipesInput) ([]Recipe, error)
	Create(ctx context.Context, value CreateRecipe) (int32, error)
	CreateIngredientEdges(ctx context.Context, value []IngredientEdge) error
}
