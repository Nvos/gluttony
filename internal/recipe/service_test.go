package recipe_test

import (
	"connectrpc.com/connect"
	"context"
	"gluttony/internal/database/sqldb"
	recipev1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/recipe"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("Create recipe", func(t *testing.T) {
		pool := sqldb.NewTestPGXPool(t)
		client := recipe.NewTestEnv(t, pool)

		input := &recipev1.CreateRecipeRequest{
			Name:          "Curry",
			Description:   "Fancy text",
			Content:       "<span>Fancy text</span>",
			Locale:        "en",
			IngredientIds: []int32{1},
		}

		_, err := client.CreateRecipe(context.Background(), connect.NewRequest(input))
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
			return
		}
	})
}
