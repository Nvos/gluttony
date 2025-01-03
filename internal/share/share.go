package share

import "strings"

type Context struct {
	IsAuthenticated bool
	Path            string
	User            *UserContext
}

func (c *Context) IsURLActive(value string) bool {
	return strings.HasPrefix(c.Path, value)
}

type UserContext struct {
	IsAdmin  bool
	Username string
	UserID   int64
}
