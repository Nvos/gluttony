package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gluttony/internal/user"
	"gluttony/pkg/password"
	"gluttony/pkg/session"
)

type Service struct {
	sessionService *session.Service
	store          user.Store
}

func NewService(store user.Store, sessionService *session.Service) *Service {
	return &Service{
		store:          store,
		sessionService: sessionService,
	}
}

func (s *Service) Create(ctx context.Context, input user.CreateInput) error {
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

func (s *Service) GetByCredentials(ctx context.Context, input user.Credentials) (user.User, error) {
	value, err := s.store.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, user.ErrInvalidCredentials
		}

		return user.User{}, fmt.Errorf("get user: %w", err)
	}

	ok, err := password.Compare(value.Password, input.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("compare password: %w", err)
	}

	if !ok {
		return user.User{}, user.ErrInvalidCredentials
	}

	return value, nil
}
