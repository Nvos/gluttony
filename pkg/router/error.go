package router

import (
	"fmt"
	"net/http"
)

type CodeError struct {
	Code int
	Err  error
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("http status code %q: %v", http.StatusText(e.Code), e.Err)
}
