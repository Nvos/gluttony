package logx

import "log/slog"

type Config struct {
	Level slog.Level
	Path  string
}
