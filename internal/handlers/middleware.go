package handlers

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/i18n"
	usersvc "gluttony/internal/service/user"
	"gluttony/internal/user"
	"gluttony/web"
	"gluttony/web/component"
	"gluttony/x/httpx"
	"gluttony/x/session"
	"log/slog"
	"net/http"
)

// ImpersonateMiddleware creates a middleware that impersonates a user by
// setting their session if not already established.
// INTENDED ONLY FOR DEV.
func ImpersonateMiddleware(
	impersonate string,
	userService *usersvc.Service,
	sessionService *session.Service,
) httpx.Middleware {
	u, err := userService.GetByUsername(context.Background(), impersonate)
	if err != nil {
		panic(err)
	}

	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			var cookie *http.Cookie
			cookie, err = c.Request.Cookie("GluttonySession")
			if err != nil {
				sess, sessErr := sessionService.New(c.Context())
				if sessErr != nil {
					return fmt.Errorf("new session: %w", sessErr)
				}

				sess.Data[user.DoerSessionKey] = u

				c.Response.Header().Add("Vary", "Cookie")
				c.Response.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
				cookie = sess.ToCookie(nil)
				c.SetCookie(cookie)
			}

			sess, err := sessionService.Get(c.Context(), cookie.Value)
			if err == nil {
				return next(c)
			}

			_, ok := user.GetSessionDoer(sess)
			if ok {
				return next(c)
			}

			ses, err := sessionService.Restore(c.Context(), cookie.Value)
			if err != nil {
				return fmt.Errorf("new session: %w", err)
			}

			ses.Data[user.DoerSessionKey] = u

			c.Response.Header().Add("Vary", "Cookie")
			c.Response.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
			c.SetCookie(ses.ToCookie(nil))

			return next(c)
		}
	}
}

func AuthenticationMiddleware(sessionService *session.Service) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
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

func AuthorizationMiddleware(role user.Role) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			doer := GetDoer(c)
			if doer == nil {
				url := "/login?next=" + c.Request.URL.Path
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
