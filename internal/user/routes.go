package user

import (
	"gluttony/x/httpx"
	"net/http"
)

func Routes(deps *Deps, mux *http.ServeMux, middlewares ...httpx.MiddlewareFunc) {
	mux.HandleFunc("GET /login", httpx.Apply(LoginViewHandler(deps), middlewares...))
	mux.HandleFunc("POST /login/form", httpx.Apply(LoginHTMXFormHandler(deps), middlewares...))
}
