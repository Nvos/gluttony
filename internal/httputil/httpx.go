package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type Problem struct {
	Status int
	Title  string
	Detail string

	err error
}

func (r *Problem) Error() string {
	return fmt.Sprintf("Problem status=%d, title=%s, detail=%s", r.Status, r.Title, r.Detail)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

func NewErrorMiddleware(logger *slog.Logger) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			routeErr := next(w, r)
			if routeErr == nil {
				return nil
			}

			var problem *Problem
			if !errors.As(routeErr, &problem) {
				problem = &Problem{
					Status: http.StatusInternalServerError,
					Title:  http.StatusText(http.StatusInternalServerError),
					err:    routeErr,
				}
			}

			w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
			w.WriteHeader(problem.Status)

			problemJSON, err := json.Marshal(problem)
			if err != nil {
				panic(fmt.Sprintf("marshalling problem json: %s", err))
			}
			_, _ = w.Write(problemJSON)

			if problem.Status != http.StatusInternalServerError {
				return nil
			}

			logger.Error(
				"Route handler unexpected failure",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("routeErr", routeErr.Error()),
			)

			return nil
		}
	}
}

func Apply(handler HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for i := range middlewares {
		handler = middlewares[i](handler)
	}

	return toStd(handler)
}

func toStd(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			panic(fmt.Sprintf("unhandled http request error: %v", err))
		}
	}
}
