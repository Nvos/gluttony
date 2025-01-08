package httpx

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

func ToHandlerFunc(h HandlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		var problem *Problem
		if !errors.As(err, &problem) {
			problem = &Problem{
				Status: http.StatusInternalServerError,
				Title:  http.StatusText(http.StatusInternalServerError),
				err:    err,
			}
		}

		w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
		w.WriteHeader(problem.Status)
		problemJSON, err := json.Marshal(problem)
		if err != nil {
			panic(fmt.Sprintf("marshalling problem json: %s", err))
		}
		_, _ = w.Write(problemJSON)

		if !(problem.Status == http.StatusInternalServerError) {
			return
		}

		logger.Error(
			"Route handler unexpected failure",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("err", problem.err.Error()),
		)
	}
}
