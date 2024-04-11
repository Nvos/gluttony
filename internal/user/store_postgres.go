package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/database"
	"gluttony/internal/database/transaction"
	"gluttony/internal/user/postgresql"
	"gluttony/internal/x/cryptox"
)

var _ Store = (*PostgresStore)(nil)

type PostgresStore struct {
	pool    postgresql.DBTX
	queries *postgresql.Queries
}

func (p *PostgresStore) UnderTransaction(tx transaction.Transaction) (Store, error) {
	pgxTx, err := transaction.GetPgxTx(tx)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		pool:    pgxTx,
		queries: p.queries.WithTx(pgxTx),
	}, nil
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

func (p *PostgresStore) Transaction(ctx context.Context) {

}

func (p *PostgresStore) Create(ctx context.Context, username, password string) (int32, error) {
	hash, err := cryptox.Argon2Hash(password, cryptox.NewDefaultArgon2Config())
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
