package config

import (
	"fmt"
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
	WebPort           int
	WebHost           string
}

func New() (Config, error) {
	if err := LoadEnv(); err != nil {
		return Config{}, fmt.Errorf("loading environment variables: %w", err)
	}

	logLevelRaw := os.Getenv("GLUTTONY_LOG_LEVEL")
	workDirectoryPath := os.Getenv("GLUTTONY_WORK_DIRECTORY_PATH")
	modeRaw := os.Getenv("GLUTTONY_MODE")
	webPortRaw := os.Getenv("GLUTTONY_WEB_PORT")
	webHostRaw := os.Getenv("GLUTTONY_WEB_HOST")

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(logLevelRaw)); err != nil {
		return Config{}, fmt.Errorf("parsing env=GLUTTONY_LOG_LEVEL: %w", err)
	}

	if !filepath.IsAbs(workDirectoryPath) {
		return Config{}, fmt.Errorf("work directory path must be absolute env=GLUTTONY_WORK_DIRECTORY_PATH")
	}

	mode, err := newMode(modeRaw)
	if err != nil {
		return Config{}, fmt.Errorf("parsing env=GLUTTONY_MODE: %w", err)
	}

	webPort, err := strconv.Atoi(webPortRaw)
	if err != nil {
		return Config{}, fmt.Errorf("parsing env=GLUTTONY_WEB_PORT: %w", err)
	}

	return Config{
		Mode:              mode,
		LogLevel:          logLevel,
		WorkDirectoryPath: workDirectoryPath,
		WebPort:           webPort,
		WebHost:           webHostRaw,
	}, nil
}
