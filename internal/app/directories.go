package app

import (
	"fmt"
	"github.com/spf13/afero"
	"gluttony/assets"
	"io/fs"
	"os"
	"path/filepath"
)

type Directories struct {
	Assets fs.FS
	Media  afero.Fs
	Root   afero.Fs
}

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

func NewDirectories(
	mode Mode,
	workRoot afero.Fs,
) (*Directories, error) {
	assetsFS, err := GetAssets(mode)
	if err != nil {
		return nil, fmt.Errorf("could not get assets: %w", err)
	}

	if err := workRoot.MkdirAll("media", os.ModePerm); err != nil {
		return nil, fmt.Errorf("could not create media directory: %w", err)
	}

	return &Directories{
		Media:  afero.NewBasePathFs(workRoot, "media"),
		Assets: assetsFS,
		Root:   workRoot,
	}, nil
}
