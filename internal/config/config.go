package config

import (
	"errors"
	"gluttony/x/httpx"
	"gluttony/x/logx"
	"gluttony/x/sqlx"
	"os"
	"path/filepath"

	"fmt"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Mode string

const (
	ModeProd Mode = "prod"
	ModeDev  Mode = "dev"
)

type Config struct {
	Mode        Mode
	Domain      string
	HTTP        httpx.Config
	Database    sqlx.Config
	Logger      logx.Config
	WorkDir     string `koanf:"work-dir"`
	Impersonate string
}

type ServerConfig struct {
	Host string
	Port int
}

func (c *Config) Validate() error {
	if c.Mode != "dev" && c.Mode != "prod" {
		return fmt.Errorf("invalid mode %q want dev or prod", c.Mode)
	}

	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.HTTP.Port)
	}

	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	if c.Database.Host == "" {
		return errors.New("database host is required")
	}

	if c.Database.User == "" {
		return errors.New("database user is required")
	}

	if c.Database.Name == "" {
		return errors.New("database name is required")
	}

	// Validate directories exist
	logDir := filepath.Dir(c.Logger.Path)
	if _, err := os.Stat(logDir); err != nil {
		return fmt.Errorf("log directory does not exist: %s", logDir)
	}

	if _, err := os.Stat(c.WorkDir); err != nil {
		return fmt.Errorf("work directory does not exist: %s", c.WorkDir)
	}

	if c.Mode == ModeProd {
		if c.Domain == "" {
			return errors.New("domain is required in production mode")
		}

		if c.Impersonate != "" {
			return errors.New("impersonate is not supported in production mode")
		}
	}

	return nil
}

func NewConfig(cfgPath string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(cfgPath), toml.Parser()); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var cfg *Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
