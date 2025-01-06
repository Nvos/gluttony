package httpx

import (
	"net/http"
)

func HTMXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}

func IsHTMXRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
