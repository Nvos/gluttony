package ingredient_test

import (
	"connectrpc.com/connect"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // "pgx" driver
	"gluttony/internal/i18n"
	"gluttony/internal/ingredient"
	ingredientv1 "gluttony/internal/proto/ingredient/v1"
	"gluttony/internal/proto/ingredient/v1/ingredientv1connect"
	"gluttony/internal/x/connectx"
	"gluttony/internal/x/httpx"
	"gluttony/internal/x/testing/assert"
	"gluttony/internal/x/testing/testdb"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newConnectEnv(t *testing.T, pool *pgxpool.Pool) ingredientv1connect.IngredientServiceClient {
	t.Helper()

	store := ingredient.NewStorePostgres(pool)
	service := ingredient.NewService(store)

	baseURL, handler, err := ingredient.NewConnectHandler(service)
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(baseURL, handler)

	server := httptest.NewUnstartedServer(httpx.ComposeMiddlewares(mux, i18n.LocaleInjectionMiddleware()))
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	return ingredient.NewConnectClient(server.Client(), server.URL)
}

func TestConnectServiceV1_Create(t *testing.T) {
	t.Run("Create valid ingredient", func(t *testing.T) {
		t.Parallel()
		client := newConnectEnv(t, testdb.NewTestPGXPool(t))

		// TODO, 01/05/2024: should return something, id or model
		_, err := client.Create(context.Background(), connect.NewRequest(&ingredientv1.CreateRequest{
			Name:   "Mango",
			Locale: "en",
			Unit:   ingredientv1.Unit_UNIT_WEIGHT,
		}))
		if assert.NilErr(t, err) {
			t.FailNow()
		}

		request := connect.NewRequest(&ingredientv1.SingleRequest{Id: 1})
		request.Header().Set("Accept-Language", "en")
		single, err := client.Single(context.Background(), request)
		if assert.NilErr(t, err) {
			t.FailNow()
		}

		assert.Equal(t, single.Msg.Id, 1)
		assert.Equal(t, single.Msg.Name, "Mango")
		assert.Equal(t, single.Msg.Unit, ingredientv1.Unit_UNIT_WEIGHT)
	})

	t.Run("Create ingredient unknown locale", func(t *testing.T) {
		t.Parallel()
		client := newConnectEnv(t, testdb.NewTestPGXPool(t))

		_, err := client.Create(context.Background(), connect.NewRequest(&ingredientv1.CreateRequest{
			Name:   "Apple",
			Locale: "de",
		}))

		connectErr, ok := connectx.AsConnectError(err)
		if assert.Equal(t, ok, true) {
			t.FailNow()
		}

		assert.Equal(t, connectErr.Code(), connect.CodeInvalidArgument)
	})
}

func TestConnectServiceV1_All(t *testing.T) {
	t.Run("All en ingredients", func(t *testing.T) {
		t.Parallel()

		pool := testdb.NewTestPGXPool(t)
		testdb.Seed(t, pool, "ingredients.sql")

		client := newConnectEnv(t, pool)

		input := connect.NewRequest(&ingredientv1.AllRequest{
			Offset: 0,
			Limit:  20,
			Search: "",
		})
		input.Header().Set("Accept-Language", "en")

		all, err := client.All(context.Background(), input)
		if assert.NilErr(t, err) {
			t.FailNow()
		}

		if assert.Equal(t, len(all.Msg.Ingredients), 3) {
			t.FailNow()
		}

		assert.Equal(t, all.Msg.Ingredients[0].Id, 1)
		assert.Equal(t, all.Msg.Ingredients[0].Name, "Apple")
		assert.Equal(t, all.Msg.Ingredients[1].Id, 2)
		assert.Equal(t, all.Msg.Ingredients[1].Name, "Chicken")
		assert.Equal(t, all.Msg.Ingredients[2].Id, 3)
		assert.Equal(t, all.Msg.Ingredients[2].Name, "Carrot")
	})
}
