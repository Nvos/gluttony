package httpx

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(*Context) error
type Middleware func(HandlerFunc) HandlerFunc

func (r *Router) toHTTPHandler(h HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := range middlewares {
		h = middlewares[i](h)
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := Context{
			Response: rw,
			Request:  req,
			Data:     make(map[string]any),
		}

		ctx.Data["Path"] = req.URL.Path

		if err := h(&ctx); err != nil {
			// Assumption is made that error should be consumed by some middleware
			panic(fmt.Sprintf("http handler: %v", err))
		}
	}
}

func WrapHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		h.ServeHTTP(ctx.Response, ctx.Request)
		return nil
	}
}
