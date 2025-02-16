package templates

import (
	"embed"
	"gluttony/internal/config"
	"io/fs"
	"os"
)

//go:embed *
var Embedded embed.FS

func GetTemplates(mode config.Mode) fs.FS {
	if mode == config.Prod {
		return Embedded
	}

	return os.DirFS("internal/web/templates")
}
