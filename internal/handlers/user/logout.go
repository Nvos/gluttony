package user

import (
	"gluttony/internal/user"
	"gluttony/x/httpx"
	"net/http"
)

func (r *Routes) LogoutHandler(c *httpx.Context) error {
	sess, ok := user.GetContextSession(c.Context())
	if !ok {
		return c.Error(http.StatusUnauthorized, nil)
	}

	if err := r.service.Logout(sess); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.SetCookie(user.NewInvalidateCookie())
	c.Redirect("/", http.StatusFound)

	return nil
}
