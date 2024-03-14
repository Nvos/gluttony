package auth

import "context"

type Store[T any] interface {
	Single(ctx context.Context, key string) (T, error)
	Create(ctx context.Context, key string, value T) error
	Delete(ctx context.Context, key string) error
}
