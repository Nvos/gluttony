package recipe

import (
	"context"
	"gluttony/internal/database/pagination"
)

type Store interface {
	Single(ctx context.Context, id int32) (Recipe, error)
	All(ctx context.Context, search string, pagination pagination.OffsetPagination) ([]Recipe, error)
	Create(ctx context.Context, value CreateRecipe) (int32, error)
}
