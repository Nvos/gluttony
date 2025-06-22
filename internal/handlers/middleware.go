package handlers

import (
	"context"
	"errors"
	"fmt"
	usersvc "gluttony/internal/service/user"
	"gluttony/internal/user"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
	"gluttony/web"
	"gluttony/web/component"
	"log/slog"
	"net/http"
	"time"
)

// ImpersonateMiddleware creates a middleware that impersonates a user by
// setting their session if not already established.
// INTENDED ONLY FOR DEV.
func ImpersonateMiddleware(
	impersonate string,
	userService *usersvc.Service,
	sessionService *session.Service,
) router.Middleware {
	u, err := userService.GetByUsername(context.Background(), impersonate)
	if err != nil {
		panic(err)
	}

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) error {
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
				const dayDuration = 24 * time.Hour
				cookie = sess.ToCookie(time.Now().UTC().Add(dayDuration))
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
			const dayDuration = 24 * time.Hour
			c.SetCookie(ses.ToCookie(time.Now().UTC().Add(dayDuration)))

			return next(c)
		}
	}
}

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

			var httpErr *router.HTTPError
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

func AuthorizationMiddleware(role user.Role) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) error {
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
