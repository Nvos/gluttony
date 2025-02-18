package handlers

import (
	"errors"
	"fmt"
	"gluttony/internal/html"
	"gluttony/internal/security"
	"gluttony/internal/web"
	"net/http"
)

type LoginForm struct {
	Username    string
	Password    string
	RedirectURL string
}

const loginView = "view/login"
const loginForm = "login/form"

func (r *Routes) LoginViewHandler(c *web.Context) error {
	redirectURL := "/recipes"
	next := c.Request.URL.Query().Get("next")
	if next != "" {
		redirectURL = next
	}

	c.Data["Form"] = LoginForm{
		RedirectURL: redirectURL,
	}

	return c.RenderView(loginView, http.StatusOK)
}

func (r *Routes) LoginHTMXFormHandler(c *web.Context) error {
	form := LoginForm{
		Username:    c.FormValue("username"),
		Password:    c.FormValue("password"),
		RedirectURL: c.FormValue("redirect_url"),
	}

	session, err := r.service.Login(c.Context(), form.Username, form.Password)
	if errors.Is(err, security.ErrInvalidCredentials) {
		c.Data["Form"] = form
		c.Data["LoginAlert"] = html.NewAlert(
			html.AlertError,
			"Invalid credentials",
			"Username and password do not match.",
		)

		return c.RenderViewFragment(loginView, loginForm, http.StatusOK)
	}

	if err != nil {
		return fmt.Errorf("login form: %w", err)
	}

	c.SetCookie(session.ToCookie())
	c.HTMXRedirect(form.RedirectURL)

	return nil
}
