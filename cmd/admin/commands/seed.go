package commands

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/pkg/database"
	"gluttony/seeds"
	"io/fs"
	"os"
)

func RunSeed(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("get work dir path: %v", err))
	}

	cfg, err := config.NewConfig(wd)
	if err != nil {
		panic(fmt.Sprintf("create config: %v", err))
	}

	pool, err := database.New(ctx, cfg.Database)
	if err != nil {
		panic(fmt.Sprintf("create db: %v", err))
	}
	defer pool.Close()

	files, err := seeds.Seeds.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read seed files: %w", err)
	}

	for i := range files {
		fmt.Println(fmt.Sprintf("Running seed %s", files[i].Name()))

		script, err := fs.ReadFile(seeds.Seeds, files[i].Name())
		if err != nil {
			return fmt.Errorf("read seed file: %w", err)
		}

		if _, err := pool.Exec(ctx, string(script)); err != nil {
			return fmt.Errorf("run seed: %w", err)
		}
	}

	return nil
}
