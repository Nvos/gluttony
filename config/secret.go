package config

import (
	"errors"
	"fmt"
	"gluttony/x/sqlx"
	"os"
	"path/filepath"
)

const postgresPasswordFilename = "postgres-password"

type Secret struct {
	Database sqlx.Secret
}

func NewSecret() (*Secret, error) {
	credentialsDirPath := os.Getenv("CREDENTIALS_DIRECTORY")
	if credentialsDirPath == "" {
		return nil, errors.New("env variable 'CREDENTIALS_DIRECTORY' should be set")
	}

	postgresPasswordPath := filepath.Join(credentialsDirPath, postgresPasswordFilename)
	if _, err := os.Stat(postgresPasswordPath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("'postgres-password' secret file does not exist")
		}

		return nil, fmt.Errorf("stat postgres-password file: %w", err)
	}

	password, err := os.ReadFile(postgresPasswordPath)
	if err != nil {
		return nil, fmt.Errorf("read postgres-password file: %w", err)
	}

	sec := &Secret{
		Database: sqlx.Secret{
			Password: string(password),
		},
	}

	if err := sec.Validate(); err != nil {
		return nil, fmt.Errorf("validate secret: %w", err)
	}

	return sec, nil
}

func (s *Secret) Validate() error {
	if s.Database.Password == "" {
		return errors.New("database password is required")
	}

	return nil
}
