package admin

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/internal/user"
	"gluttony/internal/user/postgres"
	"gluttony/x/password"
	"gluttony/x/sqlx"
)

func AddUser(
	ctx context.Context,
	cfg *config.Config,
	sec *config.Secret,
	username, pass string,
	role user.Role,
) error {
	pool, err := sqlx.New(ctx, cfg.Database, sec.Database)
	if err != nil {
		return fmt.Errorf("create db: %w", err)
	}
	defer pool.Close()

	hash, err := password.Hash(pass)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = postgres.NewStore(pool).Create(ctx, user.CreateInput{
		Username: username,
		Password: hash,
		Role:     role,
	})
	if err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	return nil
}
