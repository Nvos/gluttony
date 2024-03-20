package user

import (
	"context"
	"encoding/json"
	"errors"
	"gluttony/internal/database/transaction"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginInput struct {
	Username string
	Password string
}

type Session struct {
	UserID   int32
	Username string
}

type User struct {
	ID       int32
	Username string
	Password string `json:"-"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type user User
	out := user(u)
	out.Password = "[REDACTED]"

	return json.Marshal(out)
}

type Store interface {
	UnderTransaction(tx transaction.Transaction) (Store, error)
	Single(ctx context.Context, id int32) (User, error)
	SingleByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, username, password string) (int32, error)
}
