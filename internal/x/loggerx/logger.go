package loggerx

import (
	"log/slog"
	"os"
)

func NewLogger(level slog.Level) (*slog.Logger, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	return logger, nil
}
