package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	if store == nil {
		panic("store is nil")
	}

	return &Service{
		store: store,
	}
}

func (s *Service) New(ctx context.Context) (Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return Session{}, fmt.Errorf("generate session id: %w", err)
	}

	value := Session{
		id:   sessionID,
		Data: make(map[Key]any),
	}

	if err := s.store.Create(ctx, value); err != nil {
		return Session{}, fmt.Errorf("create session: %w", err)
	}

	return value, nil
}

func (s *Service) Get(ctx context.Context, sessionID string) (Session, error) {
	value, err := s.store.Get(ctx, sessionID)
	if err != nil {
		return Session{}, fmt.Errorf("get session: %w", err)
	}

	return value, nil
}

func (s *Service) Delete(ctx context.Context, session Session) error {
	if err := s.store.Delete(ctx, session.id); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func generateSessionID() (string, error) {
	const sessionIDLength = 32
	b := make([]byte, sessionIDLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
