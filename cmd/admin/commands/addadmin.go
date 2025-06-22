package commands

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/internal/user"
	"gluttony/internal/user/postgres"
	"gluttony/pkg/database"
	"gluttony/pkg/password"
	"os"
)

func AddAdmin(ctx context.Context, username, pass string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get work dir path: %w", err)
	}

	cfg, err := config.NewConfig(wd)
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	pool, err := database.New(ctx, cfg.Database)
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
		Role:     user.RoleAdmin,
	})
	if err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	return nil
}
