package web

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/security"
	"log/slog"
	"net/http"
)

type ReadOnlySessionStore interface {
	Get(ctx context.Context, sessionID string) (security.Session, error)
}

func AuthenticationMiddleware(store ReadOnlySessionStore) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			cookie, err := c.Request.Cookie("GluttonySession")
			if err != nil {
				return next(c)
			}

			session, err := store.Get(c.Context(), cookie.Value)
			if err != nil {
				return next(c)
			}

			c.Doer = &session
			c.Request = c.Request.WithContext(security.SetSession(c.Request.Context(), session))

			if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/" {
				c.Redirect("/recipes", http.StatusTemporaryRedirect)

				return nil
			}

			return next(c)
		}
	}
}

func ErrorMiddleware(logger *slog.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var errorCode *ErrorCode
			if errors.As(err, &errorCode) {
				c.Response.WriteHeader(errorCode.Code)
				if errorCode.err != nil {
					logger.Error("Http handler returned managed error", slog.String("err", err.Error()))
				}

				// TODO: handle 404 view redirect
				return nil
			}

			// TODO: handle 500 view redirect
			c.Response.WriteHeader(http.StatusInternalServerError)
			logger.Error("Http handler returned internal error", slog.String("err", err.Error()))

			return nil
		}
	}
}

func AuthorizationMiddleware(role security.Role) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if c.Doer == nil {
				url := fmt.Sprintf("/login?next=%s", c.Request.URL.Path)
				c.Redirect(url, http.StatusFound)

				return nil
			}

			hasAccess := role == c.Doer.Role
			if c.Doer.Role == security.Admin {
				hasAccess = true
			}

			if !hasAccess {
				url := fmt.Sprintf("/login?next=%s", c.Request.URL.Path)
				c.Redirect(url, http.StatusFound)

				return nil
			}

			return next(c)
		}
	}
}
