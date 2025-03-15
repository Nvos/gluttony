package handlers

import (
	"errors"
	"fmt"
	"gluttony/internal/user"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
	"log/slog"
	"net/http"
)

func AuthenticationMiddleware(sessionService *session.Service) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) error {
			cookie, err := c.Request.Cookie("GluttonySession")
			if err != nil {
				return next(c)
			}

			sess, err := sessionService.Get(c.Context(), cookie.Value)
			if err != nil {
				return next(c)
			}

			u, ok := user.GetSessionDoer(sess)
			if !ok {
				return next(c)
			}

			c.Data["User"] = u

			if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/" {
				c.Redirect("/recipes", http.StatusTemporaryRedirect)

				return nil
			}

			return next(c)
		}
	}
}

func ErrorMiddleware(logger *slog.Logger) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var errorCode *router.CodeError
			if errors.As(err, &errorCode) {
				c.Response.WriteHeader(errorCode.Code)
				if errorCode.Err != nil {
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

func AuthorizationMiddleware(role user.Role) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) error {
			doer := GetDoer(c)
			if doer == nil {
				url := fmt.Sprintf("/login?next=%s", c.Request.URL.Path)
				c.Redirect(url, http.StatusFound)

				return nil
			}

			hasAccess := role == doer.Role
			if doer.Role == user.RoleAdmin {
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
