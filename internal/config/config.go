package config

import (
	"errors"
	"fmt"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"io/fs"
	"os"
	"strings"
)

const envPrefix = "GLUTTONY_"
const envDelim = "_"
const envFile = ".env"

func LoadConfig(configFS fs.FS) (Config, error) {
	open, err := configFS.Open(envFile)
	if err != nil {
		return Config{}, fmt.Errorf("open config .env file: %w", err)
	}

	if err := LoadEnv(open); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env file from file=%s to env: %w", envFile, err)
	}

	k := koanf.New(".")
	if err := k.Load(env.Provider(envPrefix, envDelim, transformEnvVar), nil); err != nil {
		return Config{}, fmt.Errorf("load config from env: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config into struct: %w", err)
	}

	if err := ValidateConfig(cfg); err != nil {
		return Config{}, fmt.Errorf("config is invalid: %w", err)
	}

	return cfg, nil
}

func transformEnvVar(name string) string {
	return strings.Replace(strings.ToLower(strings.TrimPrefix(name, envPrefix)), "_", envDelim, -1)
}
