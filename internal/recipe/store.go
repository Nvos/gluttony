package recipe

import (
	"context"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
)

type IngredientEdge struct {
	RecipeID     int32
	IngredientID int32
}

type Store interface {
	UnderTransaction(tx transaction.Transaction) (Store, error)
	Single(ctx context.Context, locale i18n.Locale, id int32) (FullRecipe, error)
	All(ctx context.Context, input AllRecipesInput) ([]Recipe, error)
	Create(ctx context.Context, value CreateRecipe) (int32, error)
	CreateIngredientEdges(ctx context.Context, value []IngredientEdge) error
}
