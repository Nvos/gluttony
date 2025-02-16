package user

import (
	"context"
	"database/sql"
	"gluttony/internal/database"
	"gluttony/internal/security"
	"gluttony/internal/user/queries"
)

type Store struct {
	queries *queries.Queries
}

func NewStore(db database.DBTX) *Store {
	return &Store{
		queries: queries.New(db),
	}
}

func (s *Store) WithTx(tx *sql.Tx) *Store {
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

func (s *Store) Create(ctx context.Context, username string, password string) (int64, error) {
	userID, err := s.queries.CreateUser(ctx, queries.CreateUserParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}
