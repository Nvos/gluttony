package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	"gluttony/internal/config"
	"gluttony/internal/database/sqldb"
	"gluttony/migrations"
	"os"
)

func main() {
	ctx := context.Background()

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("use one of following commands: migrate")
		return
	}

	if err := run(ctx, args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.LoadConfig(os.DirFS(wd))
	if err != nil {
		return err
	}

	switch args[0] {
	case "migrate":
		cfg := sqldb.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Database: cfg.Database.Database,
			Options:  cfg.Database.Options,
		}

		if err := migrate(ctx, cfg); err != nil {
			return err
		}
	default:
		fmt.Println("expected: 'migrate' command")
		os.Exit(1)
	}

	return nil
}

func migrate(ctx context.Context, cfg sqldb.Config) error {
	pool, err := sqldb.NewPostgres(ctx, cfg)
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	ok, err := sqldb.IsDBRunning(ctx, db)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("database is not running")
	}

	if err := sqldb.Migrate(ctx, db, migrations.FS); err != nil {
		return err
	}

	fmt.Println("database migration complete")

	return nil
}
