package httpx

import (
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"net/http"
)

func AllowAllCORSMiddleware(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		AllowCredentials: true,
	})

	return middleware.Handler(h)
}

func ComposeMiddlewares(root http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	for i := range handlers {
		root = handlers[i](root)
	}

	return root
}
