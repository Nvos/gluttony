package html

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// collectViewRoutes collect all view routes in given fs, collection happens only
// for 2 levels of nesting.
func collectViewRoutes(f fs.FS) ([]string, error) {
	var matches []string

	roots, err := fs.Glob(f, "view/*")
	if err != nil {
		return nil, fmt.Errorf("glob view roots: %w", err)
	}

	for i := range roots {
		files, err := fs.ReadDir(f, roots[i])
		if err != nil {
			return nil, fmt.Errorf("read view files: %w", err)
		}

		isSubDir := false
		for i2 := range files {
			isSubDir = files[i2].IsDir()
			if isSubDir {
				break
			}
		}
		if !isSubDir {
			matches = append(matches, roots[i])
			continue
		}

		views, err := fs.Glob(f, fmt.Sprintf("%s/*", roots[i]))
		if err != nil {
			return nil, fmt.Errorf("glob views: %w", err)
		}

		for i2 := range views {
			if strings.HasPrefix(filepath.Base(views[i2]), "_") {
				continue
			}

			matches = append(matches, views[i2])
		}
	}

	return matches, nil
}

func parseTemplate(f fs.FS, path string) (*Template, error) {
	globs := []string{
		"base/*.gohtml",
		fmt.Sprintf("%s/*.gohtml", path),
	}

	const partCount = 3
	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) == partCount {
		// Check if any fragment directories exist
		files, err := fs.ReadDir(f, filepath.Join(parts[0], parts[1]))
		if err != nil {
			return nil, fmt.Errorf("read view files: %w", err)
		}

		for i := range files {
			if !files[i].IsDir() || !strings.HasPrefix(files[i].Name(), "_") {
				continue
			}

			glob := filepath.Join(parts[0], parts[1], files[i].Name())
			globs = append(globs, fmt.Sprintf("%s/*.gohtml", glob))
		}
	}

	out := &Template{
		name: path,
		glob: globs,
		ref:  atomic.Pointer[template.Template]{},
	}
	if err := out.Parse(f); err != nil {
		return nil, err
	}

	return out, nil
}
