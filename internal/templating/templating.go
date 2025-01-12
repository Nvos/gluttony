package templating

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"os"
	"time"
)

type Templating struct {
	baseTemplate *template.Template
	template     fs.FS
}

func New(fsys fs.FS) *Templating {
	baseTemplates := os.DirFS("internal/templating/templates")

	baseTemplate := template.New("")
	baseTemplate = baseTemplate.Funcs(template.FuncMap{
		"formatDuration": func(duration time.Duration) string {
			return time.Unix(0, 0).UTC().Add(duration).Format("15:04")
		},
		"rawHTML": func(raw string) template.HTML {
			return template.HTML(raw)
		},
		"isURL": func(s string) bool {
			_, err := url.ParseRequestURI(s)
			return err == nil
		},
	})
	baseTemplate, err := baseTemplate.ParseFS(baseTemplates, "*.fragment.gohtml")
	if err != nil {
		// TODO: handle err
		panic(fmt.Sprintf("error parsing templates: %v", err))
	}

	return &Templating{
		template:     fsys,
		baseTemplate: baseTemplate,
	}
}

// TODO: optimize for prod (need to embed, and parse at startup)
func (t *Templating) View(w io.Writer, name string, data interface{}) error {
	out, err := t.baseTemplate.Clone()
	if err != nil {
		return fmt.Errorf("cloning templates: %v", err)
	}

	glob, err := fs.Glob(t.template, "*.fragment.gohtml")
	if err != nil {
		return err
	}

	if len(glob) > 0 {
		out, err = out.ParseFS(t.template, "*.fragment.gohtml")
		if err != nil {
			return fmt.Errorf("parse template: %w", err)
		}
	}

	viewName := fmt.Sprintf("%s.view.gohtml", name)
	out, err = out.ParseFS(t.template, viewName)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	return out.ExecuteTemplate(w, viewName, data)
}

func (t *Templating) Fragment(w io.Writer, name string, data interface{}) error {
	out, err := t.baseTemplate.Clone()
	if err != nil {
		return fmt.Errorf("cloning templates: %v", err)
	}

	glob, err := fs.Glob(t.template, "*.fragment.gohtml")
	if err != nil {
		return err
	}

	if len(glob) > 0 {
		out, err = out.ParseFS(t.template, "*.fragment.gohtml")
		if err != nil {
			return fmt.Errorf("parse template: %w", err)
		}
	}

	glob, err = fs.Glob(t.template, name)
	if err != nil {
		return err
	}

	return out.ExecuteTemplate(w, name, data)
}
