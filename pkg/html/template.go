package html

import (
	"fmt"
	"html/template"
	"io/fs"
	"sync/atomic"
)

type Template struct {
	name string
	glob []string
	ref  atomic.Pointer[template.Template]
}

func (t *Template) Parse(f fs.FS) error {
	tmpl, err := template.New(t.name).
		Funcs(FuncMap()).
		ParseFS(
			f, t.glob...,
		)

	if err != nil {
		return fmt.Errorf("parsing template %s: %w", t.name, err)
	}

	t.ref.Store(tmpl)

	return nil
}
