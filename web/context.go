package web

import (
	"gluttony/user"
	"net/http"
)

type Context struct {
	Lang string
	User *user.User
	Req  *http.Request
}

func (c *Context) IsAuthenticated() bool {
	return c.User != nil
}

func NewContext(req *http.Request, user *user.User, lang string) *Context {
	return &Context{User: user, Req: req, Lang: lang}
}
