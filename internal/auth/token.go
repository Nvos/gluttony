package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random 32 byte token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
