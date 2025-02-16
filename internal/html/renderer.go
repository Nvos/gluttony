package html

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"strings"
	"time"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
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
		"hasPrefix": func(value, prefix string) bool {
			return strings.HasPrefix(value, prefix)
		},
	}
}

type TemplateName string
type FragmentName string

type Renderer struct {
	opts      RendererOptions
	fs        fs.FS
	templates map[TemplateName]*Template
}

type RendererOptions struct {
	IsReloadEnabled bool
}

func NewRenderer(f fs.FS, opts RendererOptions) (*Renderer, error) {
	matches, err := collectViewRoutes(f)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, errors.New("no view routes found")
	}

	templates := make(map[TemplateName]*Template, len(matches))
	for i := range matches {
		t, err := parseTemplate(f, matches[i])
		if err != nil {
			return nil, err
		}

		templates[TemplateName(matches[i])] = t
	}

	return &Renderer{
		opts:      opts,
		fs:        f,
		templates: templates,
	}, nil
}

func (r *Renderer) RenderTemplate(name TemplateName, w io.Writer, data any) error {
	t, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("no such template: %s", name)
	}

	if r.opts.IsReloadEnabled {
		if err := t.Parse(r.fs); err != nil {
			return fmt.Errorf("reload template %s: %v", name, err)
		}
	}

	println(t.ref.Load().DefinedTemplates())
	return t.ref.Load().ExecuteTemplate(w, "base.gohtml", data)
}

func (r *Renderer) RenderFragment(name TemplateName, fragment FragmentName, w io.Writer, data any) error {
	t, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("no such template: %s", name)
	}

	if r.opts.IsReloadEnabled {
		if err := t.Parse(r.fs); err != nil {
			return fmt.Errorf("reload template %s: %v", name, err)
		}
	}

	return t.ref.Load().Lookup(string(fragment)).Execute(w, data)
}
