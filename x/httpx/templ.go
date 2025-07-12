package httpx

import (
	"github.com/a-h/templ"
)

func (c *Context) TemplComponent(code int, component templ.Component) error {
	c.SetStatus(code)
	return component.Render(c.Context(), c.Response) //nolint:nolintlint,wrapcheck
}
