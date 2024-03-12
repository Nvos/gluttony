package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewDefaultArgon2Config() Argon2Config {
	return Argon2Config{
		Memory:      19 * 1024,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func Argon2Hash(value string, cfg Argon2Config) (string, error) {
	salt, err := generateSalt(cfg.SaltLength)
	if err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey([]byte(value), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cfg.Memory,
		cfg.Iterations,
		cfg.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func CompareArgon2(password string, encodedPassword string) (bool, error) {
	cfg, salt, hash, err := decodeArgon2Password(encodedPassword)
	if err != nil {
		return false, fmt.Errorf("decode argon2 password: %w", err)
	}

	otherHash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeArgon2Password(encodedPassword string) (Argon2Config, []byte, []byte, error) {
	parts := strings.Split(encodedPassword, "$")
	if len(parts) != 6 {
		return Argon2Config{}, nil, nil, fmt.Errorf("invalid hash")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("decode hash version: %w", err)
	}

	if version != argon2.Version {
		return Argon2Config{}, nil, nil, fmt.Errorf("incompatible argon version %d expected %d", version, argon2.Version)
	}

	cfg := &Argon2Config{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &cfg.Memory, &cfg.Iterations, &cfg.Parallelism); err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("decode hash cfg: %w", err)
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("decode hash salt: %w", err)
	}

	cfg.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("decode hash value: %w", err)
	}

	cfg.KeyLength = uint32(len(hash))

	return *cfg, salt, hash, nil
}

func generateSalt(n uint32) ([]byte, error) {
	out := make([]byte, n)
	if _, err := rand.Read(out); err != nil {
		return nil, err
	}

	return out, nil
}
