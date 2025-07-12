package user

import (
	"context"
	"errors"
	"fmt"
	"gluttony/x/session"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

const DoerSessionKey = session.Key("doer")

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func NewRole(s string) (Role, error) {
	switch s {
	case string(RoleAdmin):
		return RoleAdmin, nil
	case string(RoleUser):
		return RoleUser, nil
	default:
		return "", fmt.Errorf("invalid role = %q, must be one of: %s", s, "admin, user")
	}
}

type User struct {
	ID       int32
	Username string
	Role     Role
	Password string
}

type CreateInput struct {
	Username string
	Password string
	Role     Role
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Store interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, input CreateInput) (int32, error)
}

func GetSessionDoer(value session.Session) (User, bool) {
	return session.Get[User](value, DoerSessionKey)
}
