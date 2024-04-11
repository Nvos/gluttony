package auth

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/auth/postgresql"
	"gluttony/internal/database/sqldb"
	"gluttony/internal/database/transaction"
	"gluttony/internal/x/cryptox"
)

var _ UserStore = (*UserPostgresStore)(nil)

type UserPostgresStore struct {
	pool    postgresql.DBTX
	queries *postgresql.Queries
}

func (p *UserPostgresStore) UnderTransaction(tx transaction.Transaction) (UserStore, error) {
	pgxTx, err := transaction.GetPgxTx(tx)
	if err != nil {
		return nil, err
	}

	return &UserPostgresStore{
		pool:    pgxTx,
		queries: p.queries.WithTx(pgxTx),
	}, nil
}

func NewUserPostgresStore(pool *pgxpool.Pool) (*UserPostgresStore, error) {
	if pool == nil {
		return nil, fmt.Errorf("new postgres store: pgxpool is nil")
	}

	return &UserPostgresStore{
		pool:    pool,
		queries: postgresql.New(pool),
	}, nil
}

func (p *UserPostgresStore) Transaction(ctx context.Context) {

}

func (p *UserPostgresStore) Create(ctx context.Context, username, password string) (int32, error) {
	hash, err := cryptox.CreateHash(password)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	id, err := p.queries.CreateUser(ctx, postgresql.CreateUserParams{
		Name:     username,
		Password: hash,
	})
	if err != nil {
		return 0, fmt.Errorf("create user: %w", sqldb.TransformPgxError(err))
	}

	return id, nil
}

func (p *UserPostgresStore) Single(ctx context.Context, id int32) (User, error) {
	single, err := p.queries.SingleUser(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("single user by id=%d: %w", id, sqldb.TransformPgxError(err))
	}

	return User{
		ID:           single.ID,
		Username:     single.Name,
		PasswordHash: single.Password,
	}, nil
}

func (p *UserPostgresStore) SingleByUsername(ctx context.Context, username string) (User, error) {
	single, err := p.queries.SingleUserByName(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("single user by username=%s: %w", username, sqldb.TransformPgxError(err))
	}

	return User{
		ID:           single.ID,
		Username:     single.Name,
		PasswordHash: single.Password,
	}, nil
}
