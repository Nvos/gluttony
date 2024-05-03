package recipe_test

import (
	"connectrpc.com/connect"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/auth"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
	recipev1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/proto/recipe/v1/recipev1connect"
	"gluttony/internal/recipe"
	"gluttony/internal/x/httpx"
	"gluttony/internal/x/testing/assert"
	"gluttony/internal/x/testing/testdb"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newConnectEnv(t *testing.T, pool *pgxpool.Pool) recipev1connect.RecipeServiceClient {
	t.Helper()

	store := recipe.NewStorePostgres(pool)
	beginner := transaction.NewPgxBeginner(pool)
	service := recipe.NewService(beginner, store)

	baseURL, handler, err := recipe.NewConnectHandler(service)
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(baseURL, handler)

	server := httptest.NewUnstartedServer(
		httpx.ComposeMiddlewares(
			mux,
			i18n.LocaleInjectionMiddleware(),
			auth.TestSessionInjector(1),
		),
	)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	return recipe.NewConnectClient(server.Client(), server.URL)
}

func TestConnectService_CreateRecipe(t *testing.T) {
	t.Run("Create recipe", func(t *testing.T) {
		t.Parallel()

		pool := testdb.NewTestPGXPool(t)
		testdb.Seed(t, pool, "create_recipe_seed.sql")

		client := newConnectEnv(t, pool)

		inputCreate := &recipev1.CreateRecipeRequest{
			Name:        "Curry",
			Description: "Fancy text",
			Content:     "<span>Fancy text</span>",
			Locale:      "en",
			Ingredients: []*recipev1.CreateRecipeRequest_Ingredient{
				{
					Id:     1,
					Amount: 100,
					Count:  1,
					Note:   "small",
				},
			},
		}

		_, err := client.CreateRecipe(context.Background(), connect.NewRequest(inputCreate))
		assert.NilErr(t, err)

		// TODO, 01/05/2024:
		inputSingle := connect.NewRequest(&recipev1.SingleRecipeRequest{
			Id: 1,
		})
		inputSingle.Header().Set("Accept-Language", "en")

		got, err := client.SingleRecipe(context.Background(), inputSingle)
		if assert.NilErr(t, err) {
			t.FailNow()
		}

		assert.Equal(t, got.Msg.Id, 1)
		assert.Equal(t, got.Msg.Name, "Curry")
		assert.Equal(t, got.Msg.Description, "Fancy text")
		assert.Equal(t, got.Msg.Content, "<span>Fancy text</span>")
		assert.Equal(t, got.Msg.Ingredients[0].Id, 1)
		assert.Equal(t, got.Msg.Ingredients[0].Name, "Apple")
		assert.Equal(t, got.Msg.Ingredients[0].Amount, 100)
		assert.Equal(t, got.Msg.Ingredients[0].Count, 1)
		assert.Equal(t, got.Msg.Ingredients[0].Note, "small")
	})
}
