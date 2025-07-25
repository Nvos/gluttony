package handlers

import (
	"context"
	"errors"
	"fmt"
	"gluttony/i18n"
	user2 "gluttony/user"
	"gluttony/web"
	"gluttony/web/component"
	"gluttony/x/httpx"
	"log/slog"
	"net/http"
)

// ImpersonateMiddleware creates a middleware that impersonates a user by
// setting their session if not already established.
// INTENDED ONLY FOR DEV.
func ImpersonateMiddleware(
	impersonate string,
	userService *user2.Service,
) httpx.Middleware {
	sess, err := userService.Impersonate(context.Background(), impersonate)
	if err != nil {
		panic(err)
	}

	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			nextCtx := user2.WithContextSession(c.Request.Context(), sess)
			c.Request = c.Request.WithContext(nextCtx)

			c.Response.Header().Add("Vary", "Cookie")
			c.Response.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
			c.SetCookie(sess.ToCookie(nil))

			return next(c)
		}
	}
}

func AuthenticationMiddleware(userService *user2.Service) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			cookie, err := c.Request.Cookie(user2.SessionCookieName)
			if err != nil {
				return next(c)
			}

			sess, err := userService.GetSession(cookie.Value)
			if err != nil {
				return next(c)
			}

			nextCtx := user2.WithContextSession(c.Request.Context(), sess)
			c.Request = c.Request.WithContext(nextCtx)
			if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/" {
				c.Redirect("/recipes", http.StatusTemporaryRedirect)

				return nil
			}

			return next(c)
		}
	}
}

func ErrorMiddleware(logger *slog.Logger) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var httpErr *httpx.HTTPError
			webCtx := web.NewContext(c.Request, GetDoer(c), "en")

			if errors.As(err, &httpErr) {
				if httpErr.Err != nil {
					logger.Error(
						"Http handler",
						slog.String("code", http.StatusText(httpErr.Code)),
						slog.String("err", err.Error()),
					)
				}

				switch httpErr.Code {
				case http.StatusNotFound, http.StatusBadRequest:
					return c.TemplComponent(http.StatusNotFound, component.View404(webCtx))
				case http.StatusInternalServerError:
					return c.TemplComponent(http.StatusNotFound, component.View500(webCtx))
				default:
					c.Response.WriteHeader(httpErr.Code)

					return nil
				}
			}

			logger.Error("Http handler returned internal error", slog.String("err", err.Error()))

			return c.TemplComponent(http.StatusNotFound, component.View500(webCtx))
		}
	}
}

func AuthorizationMiddleware(role user2.Role) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			session, ok := user2.GetContextSession(c.Context())
			if !ok {
				url := "/login?next=" + c.Request.URL.Path
				c.Redirect(url, http.StatusFound)

				return nil
			}

			hasAccess := role == session.User.Role
			if session.User.Role == user2.RoleAdmin {
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

func I18nMiddleware(manager *i18n.I18n) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			// TODO: Add user setting table containing language (default should be en)
			nextCtx := i18n.WithI18nBundle(c.Context(), manager.Bundles["en"])
			nextRequest := c.Request.WithContext(nextCtx)
			c.Request = nextRequest

			return next(c)
		}
	}
}
