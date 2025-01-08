package recipe

import (
	"github.com/go-chi/chi/v5"
	"gluttony/x/httpx"
)

func Routes(deps *Deps) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/recipes/create", httpx.ToHandlerFunc(CreateViewHandler(deps), deps.logger))
		r.Post("/recipes/create/form", httpx.ToHandlerFunc(CreateFormHandler(deps), deps.logger))
		r.Get("/recipes", httpx.ToHandlerFunc(RecipesViewHandler(deps), deps.logger))
		r.Get("/recipes/{recipe_id}", httpx.ToHandlerFunc(ViewHandler(deps), deps.logger))
	}
}
