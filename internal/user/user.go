package user

import (
	"context"
	"gluttony/internal/security"
	"log/slog"
)

type User struct {
	ID       int32
	Username string
	Role     security.Role
	Password string
}

type Deps struct {
	service      *Service
	sessionStore SessionStore
	logger       *slog.Logger
}

type SessionStore interface {
	Get(ctx context.Context, key string) (security.Session, error)
	New(ctx context.Context) (security.Session, error)
	Save(ctx context.Context, value security.Session) error
	Delete(ctx context.Context, value security.Session) error
}
