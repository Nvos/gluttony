package recipe

import (
	"gluttony/internal/httputil"
	"net/http"
)

func Routes(deps *Deps, mux *http.ServeMux, middlewares ...httputil.MiddlewareFunc) {
	mux.HandleFunc("GET /recipes/create", httputil.Apply(CreateViewHandler(deps), middlewares...))
	mux.HandleFunc("POST /recipes/create/form", httputil.Apply(CreateFormHandler(deps), middlewares...))
	mux.HandleFunc("GET /recipes", httputil.Apply(RecipesViewHandler(deps), middlewares...))
	mux.HandleFunc("GET /recipes/{recipe_id}", httputil.Apply(ViewHandler(deps), middlewares...))
	mux.HandleFunc("GET /recipes/{recipe_id}/edit", httputil.Apply(EditViewHandler(deps), middlewares...))
	mux.HandleFunc("POST /recipes/edit/form", httputil.Apply(UpdateFormHandler(deps), middlewares...))
}
