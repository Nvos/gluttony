package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func LoadEnv(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	currentEnv := map[string]bool{}
	envs := os.Environ()
	for i := range envs {
		currentEnv[strings.Split(envs[i], "=")[0]] = true
	}

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		if len(parts) != 2 {
			continue
		}

		key := parts[0]

		if currentEnv[key] {
			continue
		}

		if err := os.Setenv(key, parts[1]); err != nil {
			return fmt.Errorf("set environment variable=%s: %w", key, err)
		}
	}

	return nil
}
