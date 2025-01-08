package user

import (
	"github.com/go-chi/chi/v5"
	"gluttony/x/httpx"
)

func Routes(deps *Deps) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/login", httpx.ToHandlerFunc(LoginViewHandler(deps), deps.logger))
		r.Post("/login/form", httpx.ToHandlerFunc(LoginHTMXFormHandler(deps), deps.logger))
	}
}
