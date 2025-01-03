package security

import (
	"fmt"
	"github.com/alexedwards/argon2id"
)

var defaultArgonParams = argon2id.DefaultParams

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func ComparePassword(hash, password string) error {
	ok, err := argon2id.ComparePasswordAndHash(hash, password)
	if err != nil {
		return fmt.Errorf("argon2id compare password: %w", err)
	}

	if !ok {
		return ErrInvalidCredentials
	}

	return nil
}
