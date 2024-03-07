package database

import "fmt"

type Config struct {
	Host     string
	Username string
	Password string
	Database string
}

func postgresURL(cfg Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?search_path=public&sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Database,
	)
}
