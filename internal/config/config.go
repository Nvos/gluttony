package config

import (
	"bufio"
	"errors"
	"fmt"
	"gluttony/pkg/database"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const envPrefix = "GLUTTONY_"

var defaultPaths = struct {
	logFile string
	workDir string
}{
	logFile: "/var/log/gluttony/gluttony.log",
	workDir: "/var/lib/gluttony",
}

type Environment string

const (
	EnvProduction  Environment = "prod"
	EnvDevelopment Environment = "dev"
)

type Config struct {
	Environment       Environment
	Domain            string
	Server            ServerConfig
	Database          database.Config
	Log               LogConfig
	WorkDirectoryPath string
	Impersonate       string
}

type ServerConfig struct {
	Host string
	Port int
}

type LogConfig struct {
	Level slog.Level
	Path  string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Load environment
	envMode := getEnvOrDefault("MODE", "prod")
	switch envMode {
	case string(EnvProduction):
		cfg.Environment = EnvProduction
	case string(EnvDevelopment):
		cfg.Environment = EnvDevelopment
	default:
		return nil, fmt.Errorf("invalid environment mode: %s", envMode)
	}

	// Load logging configuration
	var level slog.Level
	if err := level.UnmarshalText([]byte(getEnvOrDefault("LOG_LEVEL", "warn"))); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	cfg.Log = LogConfig{
		Level: level,
		Path:  filepath.FromSlash(getEnvOrDefault("LOG_FILE_PATH", defaultPaths.logFile)),
	}

	// Load work directory path
	cfg.WorkDirectoryPath = filepath.FromSlash(getEnvOrDefault("WORK_DIRECTORY_PATH", defaultPaths.workDir))

	// Load server config
	serverPort, err := strconv.Atoi(getEnvOrDefault("WEB_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid server port: %w", err)
	}

	cfg.Server = ServerConfig{
		Host: getEnvOrDefault("WEB_HOST", "localhost"),
		Port: serverPort,
	}

	// Load database config
	dbPort, err := strconv.Atoi(getEnvOrDefault("DATABASE_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid database port: %w", err)
	}

	cfg.Database = database.Config{
		Host:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnvOrDefault("DATABASE_USER", ""),
		Password: getEnvOrDefault("DATABASE_PASSWORD", ""),
		Name:     getEnvOrDefault("DATABASE_NAME", "gluttony"),
	}

	// Load Impersonate config only for development mode
	if cfg.Environment == EnvDevelopment {
		cfg.Impersonate = getEnvOrDefault("IMPERSONATE", "")
	}

	cfg.Domain = getEnvOrDefault("DOMAIN", "")

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}

func LoadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = only
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env file syntax at line %d: %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`) // Remove quotes if present

		if os.Getenv(key) != "" {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set cfg key %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	return nil
}

// getEnvOrDefault gets an environment variable or returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	fullKey := envPrefix + key
	if value := os.Getenv(fullKey); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	if c.Database.Host == "" {
		return errors.New("database host is required") //nolint:perfsprint
	}

	if c.Database.User == "" {
		return errors.New("database user is required")
	}

	if c.Database.Name == "" {
		return errors.New("database name is required")
	}

	// Validate directories exist
	logDir := filepath.Dir(c.Log.Path)
	if _, err := os.Stat(logDir); err != nil {
		return fmt.Errorf("log directory does not exist: %s", logDir)
	}

	if _, err := os.Stat(c.WorkDirectoryPath); err != nil {
		return fmt.Errorf("work directory does not exist: %s", c.WorkDirectoryPath)
	}

	if c.Environment == EnvProduction {
		if c.Domain == "" {
			return errors.New("domain is required in production mode")
		}

		if c.Impersonate != "" {
			return errors.New("impersonate is not supported in production mode")
		}
	}

	return nil
}
