package auth

import "context"

type SessionStore interface {
	Single(ctx context.Context, key string) (Session, error)
	Create(ctx context.Context, key string, value Session) error
	Delete(ctx context.Context, key string) error
}
