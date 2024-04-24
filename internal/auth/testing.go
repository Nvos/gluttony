package auth

import (
	"context"
	"net/http"
)

func TestSessionInjector(userID int32) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			session := Session{
				UserID: userID,
			}

			ctx := context.WithValue(request.Context(), sessionValue{}, session)
			handler.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
