package user

import (
	"gluttony/internal/httputil"
	"net/http"
)

func Routes(deps *Deps, mux *http.ServeMux, middlewares ...httputil.MiddlewareFunc) {
	mux.HandleFunc("GET /login", httputil.Apply(LoginViewHandler(deps), middlewares...))
	mux.HandleFunc("POST /login/form", httputil.Apply(LoginHTMXFormHandler(deps), middlewares...))
}
