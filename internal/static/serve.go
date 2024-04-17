package static

import (
	"io/fs"
	"net/http"
)

func FileServeHandler(path string, fs fs.FS) http.Handler {
	return http.StripPrefix(path, http.FileServerFS(fs))
}
