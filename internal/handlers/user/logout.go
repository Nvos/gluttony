package user

import (
	"gluttony/pkg/router"
	"gluttony/pkg/session"
	"net/http"
)

func (r *Routes) LogoutHandler(c *router.Context) error {
	c.SetCookie(session.NewInvalidateCookie())
	c.Redirect("/", http.StatusFound)

	return nil
}
