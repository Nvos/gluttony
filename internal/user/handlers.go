package user

import (
	"errors"
	"fmt"
	"gluttony/internal/security"
	"gluttony/internal/templating"
	"gluttony/x/httpx"
	"net/http"
)

type LoginView struct {
	Form       LoginForm
	LoginAlert templating.Alert
}

type LoginForm struct {
	Username string
	Password string
}

func LoginViewHandler(deps *Deps) httpx.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return deps.templates.View(w, "login", LoginView{})
	}
}

func LoginHTMXFormHandler(deps *Deps) httpx.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		form := LoginForm{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		session, err := deps.service.Login(r.Context(), form.Username, form.Password)
		if errors.Is(err, security.ErrInvalidCredentials) {
			model := LoginView{
				Form: form,
				LoginAlert: templating.NewAlert(
					templating.AlertError,
					"Invalid credentials",
					"Username and password do not match.",
				),
			}

			return deps.templates.Fragment(w, "login/form", model)
		}

		if err != nil {
			return fmt.Errorf("login form: %w", err)
		}

		http.SetCookie(w, session.ToCookie())
		httpx.HTMXRedirect(w, "/")

		return nil
	}
}
