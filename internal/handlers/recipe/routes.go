package recipe

import (
	"github.com/yuin/goldmark"
	"gluttony/internal/handlers"
	"gluttony/internal/service/recipe"
	"gluttony/internal/user"
	"gluttony/pkg/router"
)

type Routes struct {
	service  *recipe.Service
	markdown goldmark.Markdown
}

func NewRoutes(service *recipe.Service) (*Routes, error) {
	md := goldmark.New()

	return &Routes{
		service:  service,
		markdown: md,
	}, nil
}

func (r *Routes) Mount(mux *router.Router) {
	middlewares := []router.Middleware{
		handlers.AuthorizationMiddleware(user.RoleUser),
	}

	mux.Get("/recipes/create", r.CreateViewHandler, middlewares...)
	mux.Post("/recipes/create/form", r.CreateFormHandler, middlewares...)
	mux.Get("/recipes/{recipe_id}", r.DetailsViewHandler, middlewares...)
	mux.Get("/recipes", r.RecipesHandler, middlewares...)
	mux.Get("/recipes/{recipe_id}/update", r.UpdateViewHandler, middlewares...)
	mux.Post("/recipes/{recipe_id}/update/form", r.UpdateFormHandler, middlewares...)
}
