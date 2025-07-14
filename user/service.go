package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gluttony/config"
	"gluttony/x/password"
)

type Service struct {
	cfg            *config.Config
	sessionService *SessionService
	store          Store
}

func NewService(cfg *config.Config, store Store, sessionService *SessionService) *Service {
	return &Service{
		cfg:            cfg,
		store:          store,
		sessionService: sessionService,
	}
}

func (s *Service) Login(ctx context.Context, input Credentials) (Session, error) {
	value, err := s.GetByCredentials(ctx, input)
	if err != nil {
		return Session{}, fmt.Errorf("get user: %w", err)
	}

	session, err := s.sessionService.Create(value)
	if err != nil {
		return Session{}, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

func (s *Service) GetSession(sessionID string) (Session, error) {
	session, ok := s.sessionService.Get(sessionID)
	if !ok {
		return Session{}, fmt.Errorf("session not found")
	}

	return session, nil
}

func (s *Service) Logout(session Session) error {
	if err := s.sessionService.Delete(session.id); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) error {
	passwordHash, err := password.Hash(input.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	input.Password = passwordHash
	if _, err := s.store.Create(ctx, input); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Service) GetByCredentials(ctx context.Context, input Credentials) (User, error) {
	// Validate input
	if input.Username == "" || input.Password == "" {
		return User{}, ErrInvalidCredentials
	}

	// Check context before proceeding
	if err := ctx.Err(); err != nil {
		return User{}, fmt.Errorf("context error: %w", err)
	}

	value, err := s.store.GetByUsername(ctx, input.Username)
	if err != nil {
		// Use constant time comparison even for an error case
		_, _ = password.Compare("dummy-hash", input.Password)
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, fmt.Errorf("get user: %w", err)
	}

	ok, err := password.Compare(value.Password, input.Password)
	if err != nil {
		return User{}, fmt.Errorf("compare password: %w", err)
	}
	if !ok {
		return User{}, ErrInvalidCredentials
	}

	// Clear sensitive data before returning
	value.Password = ""
	return value, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (User, error) {
	value, err := s.store.GetByUsername(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return value, nil
}

func (s *Service) Impersonate(ctx context.Context, username string) (Session, error) {
	u, err := s.GetByUsername(ctx, username)
	if err != nil {
		return Session{}, fmt.Errorf("get session: %w", err)
	}

	session, err := s.sessionService.Create(u)
	if err != nil {
		return Session{}, fmt.Errorf("create session: %w", err)
	}

	return session, err
}
