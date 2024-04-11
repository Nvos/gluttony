package auth

import (
	"context"
	"fmt"
	"gluttony/internal/database/sqldb"
	"gluttony/internal/x/cryptox"
)

type Service struct {
	store   UserStore
	session *SessionManager
}

func NewService(store UserStore, session *SessionManager) (*Service, error) {
	if store == nil {
		return nil, fmt.Errorf("store is nil")
	}

	if session == nil {
		return nil, fmt.Errorf("session manager is nil")
	}

	return &Service{store: store, session: session}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	err := s.session.Delete(ctx, token)
	if err != nil {
		return fmt.Errorf("delete session")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.store.SingleByUsername(ctx, input.Username)
	if err != nil && sqldb.IsNotFound(err) {
		return "", ErrInvalidCredentials
	}

	if err != nil {
		return "", err
	}

	ok, err := cryptox.ComparePasswordAndHash(input.Password, user.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("compare password: %w", err)
	}

	if !ok {
		return "", ErrInvalidCredentials
	}

	session := Session{
		UserID:   user.ID,
		Username: user.Username,
	}

	token, err := s.session.Create(ctx, session)
	if err != nil {
		return "", fmt.Errorf("login: create session: %w", err)
	}

	return token, nil
}
