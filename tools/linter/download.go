package linter

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
)

func writeZipFile(dest io.Writer, zf *zip.File) error {
	f, err := zf.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	//nolint:gosec // G107 is not relevant due to this being controlled download from reliable source
	if _, err = io.Copy(dest, f); err != nil {
		return err
	}

	return nil
}

func download(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
