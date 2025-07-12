package httpx

import (
	"fmt"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(rw, req)
}

func (r *Router) Use(middlewares ...Middleware) {
	r.middlewares = append(middlewares, r.middlewares...)
}

func (r *Router) Post(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.method("POST", pattern, handler, middlewares...)
}

func (r *Router) Get(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.method("GET", pattern, handler, middlewares...)
}

func (r *Router) method(method string, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	nextMiddlewares := middlewares
	nextMiddlewares = append(nextMiddlewares, r.middlewares...)

	r.mux.HandleFunc(
		fmt.Sprintf("%v %v", method, pattern),
		r.toHTTPHandler(handler, nextMiddlewares...),
	)
}
