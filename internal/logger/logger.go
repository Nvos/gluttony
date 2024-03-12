package logger

import (
	"log/slog"
	"os"
)

func NewLogger() (*slog.Logger, slog.Leveler, error) {
	level := &slog.LevelVar{}
	level.Set(slog.LevelDebug)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	return logger, level, nil
}
