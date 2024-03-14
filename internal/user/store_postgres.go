package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/database"
	"gluttony/internal/user/postgresql"
	"gluttony/internal/util/passwordutil"
)

var _ Store = (*PostgresStore)(nil)

type PostgresStore struct {
	pool    *pgxpool.Pool
	queries *postgresql.Queries
}

func NewPostgresStore(pool *pgxpool.Pool) (*PostgresStore, error) {
	if pool == nil {
		return nil, fmt.Errorf("new postgres store: pgxpool is nil")
	}

	return &PostgresStore{
		pool:    pool,
		queries: postgresql.New(pool),
	}, nil
}

func (p *PostgresStore) Create(ctx context.Context, username, password string) (int32, error) {
	hash, err := passwordutil.Argon2Hash(password, passwordutil.NewDefaultArgon2Config())
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	id, err := p.queries.CreateUser(ctx, postgresql.CreateUserParams{
		Name:     username,
		Password: hash,
	})
	if err != nil {
		return 0, fmt.Errorf("store create user: %w", err)
	}

	return id, nil
}

func (p *PostgresStore) Single(ctx context.Context, id int32) (User, error) {
	single, err := p.queries.SingleUser(ctx, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return User{}, fmt.Errorf("postgres: user not found id=%d: %w", id, database.ErrNotFound)
	}

	if err != nil {
		return User{}, fmt.Errorf("postgres: user by id=%d: %w", id, err)
	}

	return User{
		ID:       single.ID,
		Username: single.Name,
		Password: single.Password,
	}, nil
}

func (p *PostgresStore) SingleByUsername(ctx context.Context, username string) (User, error) {
	single, err := p.queries.SingleUserByName(ctx, username)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return User{}, fmt.Errorf("postgres: user not found username=%s: %w", username, database.ErrNotFound)
	}

	if err != nil {
		return User{}, fmt.Errorf("postgres: user by username=%s: %w", username, err)
	}

	return User{
		ID:       single.ID,
		Username: single.Name,
		Password: single.Password,
	}, nil
}
