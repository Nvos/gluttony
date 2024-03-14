package auth

import (
	"context"
	"fmt"
)

type SessionManager[T any] struct {
	store Store[T]
}

func (s *SessionManager[T]) Single(ctx context.Context, key string) (T, error) {
	single, err := s.store.Single(ctx, key)
	if err != nil {
		var t T

		return t, fmt.Errorf("single session by key=%s: %w", key, err)
	}

	return single, nil
}

func (s *SessionManager[T]) Delete(ctx context.Context, key string) error {
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (s *SessionManager[T]) Create(ctx context.Context, value T) (string, error) {
	token, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}

	if err := s.store.Create(ctx, token, value); err != nil {
		return "", fmt.Errorf("set session value: %w", err)
	}

	return token, nil
}

func NewSessionManager[T any](store Store[T]) *SessionManager[T] {
	return &SessionManager[T]{
		store: store,
	}
}
