package user

import "github.com/go-chi/chi/v5"

func Routes(deps *Deps) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/login", LoginViewHandler(deps))
		r.Post("/login/form", LoginHTMXFormHandler(deps))
	}
}
