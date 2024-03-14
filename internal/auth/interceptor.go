package auth

import (
	"context"
	"fmt"
	"net/http"
)

type sessionValue struct{}
type sessionToken struct{}

func SessionHttpMiddleware[T any](sm *SessionManager[T]) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			cookie, err := request.Cookie(cookieName)
			if err == nil {
				single, err := sm.Single(request.Context(), cookie.Value)
				if err == nil {
					nextCtx := context.WithValue(request.Context(), sessionToken{}, cookie.Value)
					nextCtx = context.WithValue(nextCtx, sessionValue{}, single)
					handler.ServeHTTP(writer, request.WithContext(nextCtx))
					return
				}
			}

			handler.ServeHTTP(writer, request)
		})
	}
}

func GetSession[T any](ctx context.Context) (T, error) {
	req, ok := ctx.Value(sessionValue{}).(T)
	if !ok {
		var t T
		return t, fmt.Errorf("ctx has to 'sessionValue' value")
	}

	return req, nil
}

func GetSessionToken(ctx context.Context) (string, error) {
	req, ok := ctx.Value(sessionToken{}).(string)
	if !ok {
		return "", fmt.Errorf("ctx has to 'sessionToken' value")
	}

	return req, nil
}
