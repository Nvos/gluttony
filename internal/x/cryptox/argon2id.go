package cryptox

import (
	"github.com/alexedwards/argon2id"
)

func CreateHash(value string) (string, error) {
	return argon2id.CreateHash(value, argon2id.DefaultParams)
}

func ComparePasswordAndHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
