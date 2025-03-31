package user

import (
	"errors"
	"fmt"
	"gluttony/internal/handlers"
	"gluttony/internal/user"
	"gluttony/pkg/router"
	"net/http"
	"time"
)

type LoginForm struct {
	Username    string
	Password    string
	RedirectURL string
}

func (r LoginForm) ToCredentials() user.Credentials {
	return user.Credentials{
		Username: r.Username,
		Password: r.Password,
	}
}

const loginView = "view/login"
const loginForm = "login/form"

func (r *Routes) LoginViewHandler(c *router.Context) error {
	redirectURL := "/recipes"
	next := c.Request.URL.Query().Get("next")
	if next != "" {
		redirectURL = next
	}

	c.Data["Form"] = LoginForm{
		Username:    "",
		Password:    "",
		RedirectURL: redirectURL,
	}

	return c.RenderView(loginView, http.StatusOK)
}

func (r *Routes) LoginHTMXFormHandler(c *router.Context) error {
	form := LoginForm{
		Username:    c.FormValue("username"),
		Password:    c.FormValue("password"),
		RedirectURL: c.FormValue("redirect_url"),
	}

	u, err := r.service.GetByCredentials(c.Context(), form.ToCredentials())
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			c.Data["Form"] = form
			c.Data["LoginAlert"] = handlers.NewAlert(
				handlers.AlertError,
				"Invalid credentials",
				"Username and password do not match.",
			)

			return c.RenderViewFragment(loginView, loginForm, http.StatusOK)
		}

		return fmt.Errorf("get user=%s: %w", form.Username, err)
	}

	session, err := r.sessionService.New(c.Context())
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}

	session.Data[user.DoerSessionKey] = u

	c.Response.Header().Add("Vary", "Cookie")
	c.Response.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
	const dayDuration = 24 * time.Hour
	c.SetCookie(session.ToCookie(time.Now().UTC().Add(dayDuration)))
	c.HTMXRedirect(form.RedirectURL)

	return nil
}
