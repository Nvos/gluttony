package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type SessionManager struct {
	store SessionStore
}

func (s *SessionManager) Single(ctx context.Context, key string) (Session, error) {
	single, err := s.store.Single(ctx, key)
	if err != nil {
		return Session{}, fmt.Errorf("single session by key=%s: %w", key, err)
	}

	return single, nil
}

func (s *SessionManager) Delete(ctx context.Context, key string) error {
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (s *SessionManager) Create(ctx context.Context, value Session) (string, error) {
	token, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}

	if err := s.store.Create(ctx, token, value); err != nil {
		return "", fmt.Errorf("set session value: %w", err)
	}

	return token, nil
}

func NewSessionManager(store SessionStore) *SessionManager {
	return &SessionManager{
		store: store,
	}
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random 32 byte token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
