package app

import (
	"fmt"
	"gluttony/assets"
	"io/fs"
	"os"
	"path/filepath"
)

func GetAssets(mode Mode) (fs.FS, error) {
	if mode == Prod {
		return assets.Embedded, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get executable path: %w", err)
	}

	return os.DirFS(filepath.Join(wd, "assets")), nil
}
