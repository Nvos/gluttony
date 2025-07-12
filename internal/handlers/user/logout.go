package user

import (
	"gluttony/x/httpx"
	"gluttony/x/session"
	"net/http"
)

func (r *Routes) LogoutHandler(c *httpx.Context) error {
	c.SetCookie(session.NewInvalidateCookie())
	c.Redirect("/", http.StatusFound)

	return nil
}
