package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gluttony/internal/user"
)

var _ user.Store = (*Store)(nil)

type Store struct {
	queries *Queries
}

func NewStore(db DBTX) *Store {
	return &Store{
		queries: New(db),
	}
}

func (s *Store) WithTx(tx pgx.Tx) *Store {
	return &Store{
		queries: s.queries.WithTx(tx),
	}
}

func (s *Store) GetByUsername(ctx context.Context, username string) (user.User, error) {
	value, err := s.queries.GetUser(ctx, username)
	if err != nil {
		return user.User{}, fmt.Errorf("get user by username %q: %w", username, err)
	}

	return user.User{
		ID:       value.ID,
		Username: value.Username,
		Role:     newRole(value.Role),
		Password: value.Password,
	}, nil
}

func (s *Store) Create(ctx context.Context, input user.CreateInput) (int32, error) {
	userID, err := s.queries.CreateUser(ctx, CreateUserParams{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return userID, nil
}

func newRole(role UsersRole) user.Role {
	switch role {
	case UsersRoleAdmin:
		return user.RoleAdmin
	case UsersRoleUser:
		return user.RoleUser
	}

	panic(fmt.Sprintf("unknown role %s", role))
}
