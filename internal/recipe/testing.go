package recipe

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/auth"
	"gluttony/internal/database/transaction"
	"gluttony/internal/i18n"
	"gluttony/internal/proto/recipe/v1/recipev1connect"
	"gluttony/internal/x/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		`INSERT INTO users(id, name, password) VALUES (1, 'user', 'password')`,
	)
	if err != nil {
		t.Fatalf("user seed failed: %v", err)
	}

	_, err = pool.Exec(
		context.Background(),
		`INSERT INTO ingredients(id, name) VALUES (1, '{"en": "Apple", "pl": "Jabłko"}')`,
	)
	if err != nil {
		t.Fatalf("ingredient seed failed: %v", err)
	}
}

func NewTestEnv(t *testing.T, pool *pgxpool.Pool) recipev1connect.RecipeServiceClient {
	t.Helper()

	store := NewStorePostgres(pool)
	beginner := transaction.NewPgxBeginner(pool)
	service := NewService(beginner, store)

	baseURL, handler, err := NewConnectHandler(service)
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

	seed(t, pool)

	return NewConnectClient(server.Client(), server.URL)
}
