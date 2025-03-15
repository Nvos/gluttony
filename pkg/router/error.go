package router

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code int
	Err  error
}

func WithError(err error) func(*HTTPError) {
	return func(httpErr *HTTPError) {
		httpErr.Err = err
	}
}

func NewHTTPError(code int, opts ...func(err *HTTPError)) *HTTPError {
	e := &HTTPError{Code: code}
	for i := range opts {
		opts[i](e)
	}

	return e
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http status code %q: %v", http.StatusText(e.Code), e.Err)
}
