package user

import (
	"errors"
	"fmt"
	"gluttony/internal/security"
	"gluttony/internal/templates"
	"gluttony/x/httpx"
	"net/http"
)

type LoginView struct {
	Form       LoginForm
	LoginAlert templates.Alert
}

type LoginForm struct {
	Username string
	Password string
}

func LoginViewHandler(deps *Deps) httpx.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		t, err := deps.templates.Get("user", "login")
		if err != nil {
			return fmt.Errorf("login get template: %w", err)
		}

		model := LoginView{}
		if err = t.View(w, model); err != nil {
			return fmt.Errorf("expected template login to exist: %w", err)
		}

		return nil
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
			t, err := deps.templates.Get("user", "login")
			if err != nil {
				return fmt.Errorf("expected template login to exist: %w", err)
			}

			model := LoginView{
				Form: form,
				LoginAlert: templates.NewAlert(
					templates.AlertError,
					"Invalid credentials",
					"Username and password do not match.",
				),
			}

			if err = t.Fragment(w, "login/form", model); err != nil {
				return fmt.Errorf("expected template login/form to exist: %w", err)
			}

			return nil
		}

		if err != nil {
			return fmt.Errorf("login form: %w", err)
		}

		http.SetCookie(w, session.ToCookie())
		httpx.HTMXRedirect(w, "/")

		return nil
	}
}
