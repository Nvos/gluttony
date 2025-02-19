package recipe

import (
	"github.com/yuin/goldmark"
	"gluttony/internal/recipe"
	"gluttony/internal/security"
	"gluttony/internal/web"
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

func (r *Routes) Mount(router *web.Router) {
	middlewares := []web.Middleware{
		web.AuthorizationMiddleware(security.User),
	}

	router.Get("/recipes/create", r.CreateViewHandler, middlewares...)
	router.Post("/recipes/create/form", r.CreateFormHandler, middlewares...)
	router.Get("/recipes/{recipe_id}", r.DetailsViewHandler, middlewares...)
	router.Get("/recipes", r.ListViewHandler, middlewares...)
	router.Get("/recipes/{recipe_id}/update", r.UpdateViewHandler, middlewares...)
	router.Post("/recipes/{recipe_id}/update/form", r.UpdateFormHandler, middlewares...)
}
