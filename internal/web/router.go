package web

import (
	"fmt"
	"gluttony/internal/html"
	"net/http"
)

type Router struct {
	mux      *http.ServeMux
	renderer *html.Renderer
}

func NewRouter(renderer *html.Renderer) *Router {
	return &Router{mux: http.NewServeMux(), renderer: renderer}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(rw, req)
}

func (r *Router) Post(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.mux.HandleFunc(fmt.Sprintf("POST %v", pattern), r.toHttpHandler(handler, middlewares...))
}

func (r *Router) Get(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.mux.HandleFunc(fmt.Sprintf("GET %v", pattern), r.toHttpHandler(handler, middlewares...))
}
