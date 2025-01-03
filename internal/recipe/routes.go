package recipe

import (
	"github.com/go-chi/chi/v5"
)

func Routes(deps *Deps) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/recipes/create", RecipeCreateViewHandler(deps))
		r.Post("/recipes/create/form", RecipesCreateHandler(deps))
		r.Get("/recipes", RecipesViewHandler(deps))
	}
}
