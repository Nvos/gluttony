package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"path/filepath"
	"strings"
)

const fileConfigName = "config.toml"
const envPrefix = "GLUTTONY_"
const envDelim = "."
const envFile = ".env"

func LoadConfig(path string) (Config, error) {
	configFilePath := filepath.Join(path, fileConfigName)
	envFilePath := filepath.Join(path, envFile)

	if err := godotenv.Load(envFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env file from path=%s to env: %w", path, err)
	}

	k := koanf.New(".")

	if err := k.Load(file.Provider(configFilePath), toml.Parser()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load config from file path=%s: %w", configFilePath, err)
	}

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
