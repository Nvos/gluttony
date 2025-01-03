package share

import (
	"context"
	"gluttony/internal/security"
	"net/http"
)

const contextID = "Context"

func MustGetContext(ctx context.Context) *Context {
	got, ok := ctx.Value(contextID).(*Context)
	if !ok {
		panic("expected context to exist")
	}

	return got
}

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := security.GetSession(r.Context())

		ctx := &Context{
			Path: r.URL.Path,
		}
		if ok {
			ctx.IsAuthenticated = true
			ctx.User = &UserContext{
				IsAdmin:  false,
				Username: session.Username,
				UserID:   session.UserID,
			}
		}

		nextCtx := context.WithValue(r.Context(), contextID, ctx)
		next.ServeHTTP(w, r.WithContext(nextCtx))
	})
}
