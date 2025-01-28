package recipe

import (
	"gluttony/x/httpx"
	"net/http"
)

func Routes(deps *Deps, mux *http.ServeMux, middlewares ...httpx.MiddlewareFunc) {
	mux.HandleFunc("GET /recipes/create", httpx.Apply(CreateViewHandler(deps), middlewares...))
	mux.HandleFunc("POST /recipes/create/form", httpx.Apply(CreateFormHandler(deps), middlewares...))
	mux.HandleFunc("GET /recipes", httpx.Apply(RecipesViewHandler(deps), middlewares...))
	mux.HandleFunc("GET /recipes/{recipe_id}", httpx.Apply(ViewHandler(deps), middlewares...))
}
