package user

import (
	"context"
	"fmt"
	"gluttony/internal/auth"
	"gluttony/internal/database"
	"gluttony/internal/util/passwordutil"
)

type Service struct {
	store   Store
	session *auth.SessionManager[Session]
}

func NewService(store Store, session *auth.SessionManager[Session]) (*Service, error) {
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
	if err != nil && database.IsNotFound(err) {
		return "", ErrInvalidCredentials
	}

	if err != nil {
		return "", fmt.Errorf("login: user by username: %w", err)
	}

	ok, err := passwordutil.CompareArgon2(input.Password, user.Password)
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
