package user

import (
	"context"
	"github.com/jackc/pgx/v5"
	"gluttony/internal/security"
	"gluttony/internal/user/postgres"
)

type Store struct {
	queries *postgres.Queries
}

func NewStore(db postgres.DBTX) *Store {
	return &Store{
		queries: postgres.New(db),
	}
}

func (s *Store) WithTx(tx pgx.Tx) *Store {
	return &Store{
		queries: s.queries.WithTx(tx),
	}
}

func (s *Store) GetByUsername(ctx context.Context, username string) (User, error) {
	user, err := s.queries.GetUser(ctx, username)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:       user.ID,
		Username: user.Username,
		Role:     security.Role(user.Role),
		Password: user.Password,
	}, nil
}

func (s *Store) Create(ctx context.Context, username string, password string) (int32, error) {
	userID, err := s.queries.CreateUser(ctx, postgres.CreateUserParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}
