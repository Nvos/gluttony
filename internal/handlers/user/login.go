package user

import (
	"encoding/json"
	"errors"
	"fmt"
	datastar "github.com/starfederation/datastar/sdk/go"
	"gluttony/internal/handlers"
	"gluttony/internal/user"
	"gluttony/web"
	"gluttony/web/component"
	"gluttony/x/httpx"
	"net/http"
)

func (r *Routes) LoginViewHandler(c *httpx.Context) error {
	redirectURL := "/recipes"
	next := c.Request.URL.Query().Get("next")
	if next != "" {
		redirectURL = next
	}

	webCtx := web.NewContext(c.Request, handlers.GetDoer(c), "en")

	formProps := component.LoginFormProps{
		Credentials: user.Credentials{
			Username: "",
			Password: "",
		},
		RedirectURL: redirectURL,
	}

	return c.TemplComponent(http.StatusOK, component.ViewLogin(webCtx, formProps))
}

func (r *Routes) LoginFormHandler(c *httpx.Context) error {
	var props component.LoginFormProps
	if err := json.NewDecoder(c.Request.Body).Decode(&props); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	sess, err := r.service.Login(c.Context(), props.Credentials)
	if err != nil {
		sse := datastar.NewSSE(c.Response, c.Request)

		if errors.Is(err, user.ErrInvalidCredentials) {
			alert := component.NewAlert(
				component.AlertError,
				"Invalid credentials",
				"Username and password do not match.",
			)
			props.Credentials.Password = ""

			err := sse.MergeFragmentTempl(
				component.Alert(alert),
				datastar.WithSelectorID("alert"),
			)
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}

			return nil
		}

		return fmt.Errorf("get user=%s: %w", props.Credentials.Username, err)
	}

	c.Response.Header().Add("Vary", "Cookie")
	c.Response.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
	c.SetCookie(sess.ToCookie(r.cfg))

	sse := datastar.NewSSE(c.Response, c.Request)
	if err := sse.Redirect(props.RedirectURL); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return nil
}
