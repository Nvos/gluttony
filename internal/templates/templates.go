package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"time"
)

type TemplateExecutor interface {
	View(w io.Writer, data interface{}) error
	Fragment(w io.Writer, name string, data interface{}) error
}

type Executor struct {
	template *template.Template
}

func (e *Executor) View(w io.Writer, data interface{}) error {
	println(e.template.DefinedTemplates())
	return e.template.ExecuteTemplate(w, "base.gohtml", data)
}

func (e *Executor) Fragment(w io.Writer, name string, data interface{}) error {
	println(e.template.DefinedTemplates())
	return e.template.ExecuteTemplate(w, name, data)
}

type Templates struct {
	scopes map[string]fs.FS
}

func New(scopes map[string]fs.FS) *Templates {
	return &Templates{scopes: scopes}
}

func (t *Templates) Get(scope, name string) (TemplateExecutor, error) {
	baseFS, ok := t.scopes["base"]
	if !ok {
		panic("expected base scope template fs to exist")
	}

	templ := template.New("")
	templ = templ.Funcs(template.FuncMap{
		"formatDuration": func(duration time.Duration) string {
			return time.Unix(0, 0).UTC().Add(duration).Format("15:04")
		},
		"rawHTML": func(raw string) template.HTML {
			return template.HTML(raw)
		},
	})

	templ, err := templ.ParseFS(baseFS, "*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("parse base template: %w", err)
	}

	scopeFS, ok := t.scopes[scope]
	if !ok {
		return nil, fmt.Errorf("get scope template, name = %s", scope)
	}

	templ, err = templ.ParseFS(scopeFS, name+".gohtml")
	if err != nil {
		return nil, fmt.Errorf("parse scope template, name = %s: %w", scope, err)
	}

	return &Executor{template: templ}, nil
}
