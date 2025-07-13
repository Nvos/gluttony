package slogx

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Level slog.Level
	Path  string
}

func NewDev(level slog.Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       level,
		AddSource:   false,
		ReplaceAttr: nil,
	}))
}

func NewProd(level slog.Level, filePath string) (*slog.Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open log file, path='%s': %w", filePath, err)
	}

	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:       level,
		AddSource:   false,
		ReplaceAttr: nil,
	})), nil
}
