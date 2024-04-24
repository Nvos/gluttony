package ingredient

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/i18n"
	"gluttony/internal/proto/ingredient/v1/ingredientv1connect"
	"gluttony/internal/x/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewTestEnv(t *testing.T, pool *pgxpool.Pool) ingredientv1connect.IngredientServiceClient {
	t.Helper()

	store := NewStorePostgres(pool)
	service := NewService(store)

	baseURL, handler, err := NewConnectHandler(service)
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(baseURL, handler)

	server := httptest.NewUnstartedServer(httpx.ComposeMiddlewares(mux, i18n.LocaleInjectionMiddleware(nil)))
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	return NewConnectClient(server.Client(), server.URL)
}
