package app

import (
	"errors"
	"fmt"
	"gluttony/pkg/database"
	"gluttony/pkg/env"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type Mode string

const Prod Mode = "prod"
const Dev Mode = "dev"

func newMode(value string) (Mode, error) {
	switch value {
	case "prod":
		return Prod, nil
	case "dev":
		return Dev, nil
	}

	return "", fmt.Errorf("invalid mode: %s", value)
}

type Config struct {
	Mode              Mode
	LogLevel          slog.Level
	WorkDirectoryPath string
	LogFilePath       string
	WebPort           int
	WebHost           string
	Database          database.Config
}

func (c Config) Validate() error {
	if !filepath.IsAbs(c.WorkDirectoryPath) {
		return errors.New("work directory path must be absolute env=GLUTTONY_WORK_DIRECTORY_PATH")
	}

	if !filepath.IsAbs(c.LogFilePath) {
		return errors.New("log path must be absolute env=GLUTTONY_LOG_FILE_PATH")
	}

	return nil
}

func NewConfig() (Config, error) {
	if err := env.LoadEnv(); err != nil {
		return Config{}, fmt.Errorf("loading environment variables: %w", err)
	}

	mustGet := func(name string) string {
		key := "GLUTTONY_" + name
		got := os.Getenv(key)
		if got == "" {
			panic(fmt.Sprintf("config %v is not set", key))
		}

		return got
	}

	mode, err := newMode(mustGet("MODE"))
	if err != nil {
		return Config{}, fmt.Errorf("parsing env=GLUTTONY_MODE: %w", err)
	}

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(mustGet("LOG_LEVEL"))); err != nil {
		return Config{}, fmt.Errorf("parsing env=GLUTTONY_LOG_LEVEL: %w", err)
	}

	webPort, err := strconv.Atoi(mustGet("WEB_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("parsing env=GLUTTONY_WEB_PORT: %w", err)
	}

	cfg := Config{
		Mode:              mode,
		LogLevel:          logLevel,
		WorkDirectoryPath: mustGet("WORK_DIRECTORY_PATH"),
		LogFilePath:       mustGet("LOG_FILE_PATH"),
		WebPort:           webPort,
		WebHost:           mustGet("WEB_HOST"),
		Database: database.Config{
			Name:     mustGet("DATABASE_NAME"),
			User:     mustGet("DATABASE_USER"),
			Host:     mustGet("DATABASE_HOST"),
			Port:     mustGet("DATABASE_PORT"),
			Password: mustGet("DATABASE_PASSWORD"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}
