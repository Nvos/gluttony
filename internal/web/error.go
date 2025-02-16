package web

import (
	"fmt"
	"net/http"
)

type ErrorCode struct {
	Code int
	err  error
}

func (e *ErrorCode) Error() string {
	return fmt.Sprintf("http status code %q: %w", http.StatusText(e.Code), e.err)
}
