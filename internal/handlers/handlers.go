package handlers

import (
	"gluttony/internal/user"
	"gluttony/pkg/router"
)

func GetDoer(c *router.Context) *user.User {
	value, ok := c.Data["User"]
	if !ok {
		return nil
	}

	doer, ok := value.(user.User)
	if !ok {
		panic("Invalid type set in 'User' router ctx key")
	}

	return &doer
}
