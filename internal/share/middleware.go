package share

import (
	"context"
	"gluttony/internal/httputil"
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

func ContextMiddleware(next httputil.HandlerFunc) httputil.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
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
		return next(w, r.WithContext(nextCtx))
	}
}
