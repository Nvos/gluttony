package router

import (
	"context"
	"fmt"
	"gluttony/pkg/html"
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Data     map[string]any
	renderer *html.Renderer
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}

func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

func (c *Context) Form() error {
	if err := c.Request.ParseForm(); err != nil {
		return &CodeError{http.StatusBadRequest, err}
	}

	return nil
}

func (c *Context) SetStatus(code int) {
	c.Response.WriteHeader(code)
}

func (c *Context) Error(code int, err error) error {
	return &CodeError{code, err}
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
}

func (c *Context) HTMXRedirect(url string) {
	c.Response.Header().Set("Hx-Redirect", url)
}

func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}

func (c *Context) RenderView(name html.TemplateName, code int) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(code)

	if err := c.renderer.RenderTemplate(name, c.Response, c.Data); err != nil {
		return fmt.Errorf("execute template %q: %w", name, err)
	}

	return nil
}

func (c *Context) SetData(key string, value any) {
	c.Data[key] = value
}

func (c *Context) RenderViewFragment(name html.TemplateName, fragment html.FragmentName, code int) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(code)
	if err := c.renderer.RenderFragment(name, fragment, c.Response, c.Data); err != nil {
		return fmt.Errorf("execute template %q: %w", name, err)
	}

	return nil
}

// FormString returns the first value matching the provided key in the form as a string.
func (c *Context) FormString(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) FormParse() error {
	const maxFormSize = 32 << 20
	if err := c.Request.ParseMultipartForm(maxFormSize); err != nil {
		return fmt.Errorf("parse multipart form: %w", err)
	}

	return nil
}

// FormStrings returns a string slice for the provided key from the form.
func (c *Context) FormStrings(key string) []string {
	if c.Request.Form == nil {
		panic("FormStrings called with nil request form")
	}

	if v, ok := c.Request.Form[key]; ok {
		return v
	}

	return nil
}

func (c *Context) IsHTMXRequest() bool {
	return c.Request.Header.Get("Hx-Request") == "true"
}
