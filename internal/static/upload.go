package static

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const MAX_UPLOAD_KB = 1024 * 1204

func UploadHandler(dataDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_KB)
		if err := r.ParseMultipartForm(MAX_UPLOAD_KB); err != nil {
			// TODO, 17/04/2024: log stuff
			panic(err)
		}

		// Handle paths in safe way
		file, _, err := r.FormFile("file")
		if err != nil {
			// TODO, 17/04/2024: log stuff
			panic(err)
		}
		defer file.Close()

		ok, err := isFileImage(file)
		if err != nil {
			panic(err)
		}

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		fileName := uuid.New().String()
		destination, err := os.Create(filepath.Join(dataDir, fileName))
		if err != nil {
			panic(err)
		}
		defer destination.Close()

		if _, err := io.Copy(destination, file); err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("/" + path.Join("static", fileName)))
	})
}

func isFileImage(file multipart.File) (bool, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return false, fmt.Errorf("read initial 512 bytes of upload file: %w", err)
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return false, nil
	}

	if _, err := file.Seek(0, 0); err != nil {
		return false, fmt.Errorf("return to pos 0,0 of upload file: %w", err)
	}

	return true, nil
}
