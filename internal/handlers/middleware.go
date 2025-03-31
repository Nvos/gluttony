package handlers

import (
	"errors"
	"fmt"
	"gluttony/internal/user"
	"gluttony/pkg/html"
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

const view500 html.TemplateName = "view/500"
const view404 html.TemplateName = "view/404"

func ErrorMiddleware(logger *slog.Logger) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var httpErr *router.HTTPError
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
					return c.RenderView(view404, http.StatusNotFound)
				case http.StatusInternalServerError:
					return c.RenderView(view500, http.StatusInternalServerError)
				default:
					c.Response.WriteHeader(httpErr.Code)

					return nil
				}
			}

			logger.Error("Http handler returned internal error", slog.String("err", err.Error()))
			return c.RenderView(view500, http.StatusInternalServerError)
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
