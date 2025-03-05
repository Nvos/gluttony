package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/security"
)

type Service struct {
	db      *Store
	session SessionStore
}

func NewService(db *pgxpool.Pool, sessionStore SessionStore) *Service {
	return &Service{
		db:      NewStore(db),
		session: sessionStore,
	}
}

func (s *Service) Create(ctx context.Context, username, password string) error {
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = s.db.Create(ctx, username, passwordHash)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Service) Logout(ctx context.Context) error {
	session, ok := security.GetSession(ctx)
	if !ok {
		return security.ErrSessionNotFound
	}

	if err := s.session.Delete(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, username, password string) (security.Session, error) {
	user, err := s.db.GetByUsername(ctx, username)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return security.Session{}, security.ErrInvalidCredentials
	}

	if err != nil {
		return security.Session{}, err
	}

	err = security.ComparePassword(password, user.Password)
	if errors.Is(err, security.ErrInvalidCredentials) {
		return security.Session{}, err
	}

	if err != nil {
		return security.Session{}, fmt.Errorf("compare password: %w", err)
	}

	session, err := s.session.New(ctx)
	if err != nil {
		return security.Session{}, fmt.Errorf("create session: %w", err)
	}

	session.UserID = user.ID
	session.Username = user.Username
	session.Role = user.Role

	if err := s.session.Save(ctx, session); err != nil {
		return security.Session{}, fmt.Errorf("save session: %w", err)
	}

	return session, nil
}
