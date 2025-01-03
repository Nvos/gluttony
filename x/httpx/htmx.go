package httpx

import (
	"net/http"
)

func HTMXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}
