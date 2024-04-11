package auth

import (
	"context"
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
	ID           int32
	Username     string
	PasswordHash string
}

type UserStore interface {
	UnderTransaction(tx transaction.Transaction) (UserStore, error)
	Single(ctx context.Context, id int32) (User, error)
	SingleByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, username, password string) (int32, error)
}
