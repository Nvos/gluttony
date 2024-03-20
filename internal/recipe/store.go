package recipe

import (
	"context"
	"gluttony/internal/database/pagination"
	"gluttony/internal/database/transaction"
)

type Store interface {
	UnderTransaction(tx transaction.Transaction) (Store, error)
	Single(ctx context.Context, id int32) (FullRecipe, error)
	All(ctx context.Context, search string, pagination pagination.OffsetPagination) ([]Recipe, error)
	Create(ctx context.Context, value CreateRecipe) (int32, error)
	CreateRecipeSteps(ctx context.Context, recipeID int32, steps []CreateStep) error
}
