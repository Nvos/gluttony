package recipe

import (
	"github.com/go-chi/chi/v5"
)

func Routes(deps *Deps) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/recipes/create", CreateViewHandler(deps))
		r.Post("/recipes/create/form", CreateFormHandler(deps))
		r.Get("/recipes", RecipesViewHandler(deps))
		r.Get("/recipes/{recipe_id}", ViewHandler(deps))
	}
}
