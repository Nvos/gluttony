package cryptox

import (
	"github.com/alexedwards/argon2id"
)

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func CreateHash(value string) (string, error) {
	return argon2id.CreateHash(value, argon2id.DefaultParams)
}

func ComparePasswordAndHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
