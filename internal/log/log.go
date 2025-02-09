package log

import (
	"gluttony/internal/config"
	"io"
	"log/slog"
	"os"
)

type TeeWriter struct {
	writers []io.Writer
}

func (t *TeeWriter) Write(p []byte) (n int, err error) {
	for i := range t.writers {
		n, err = t.writers[i].Write(p)
		if err != nil {
			return n, err
		}
	}

	return n, err
}

func New(
	mode config.Mode,
	level slog.Level,
	writer io.Writer,
) *slog.Logger {
	if mode == config.Dev {
		tee := &TeeWriter{
			writers: []io.Writer{writer, os.Stdout},
		}

		return slog.New(slog.NewTextHandler(tee, &slog.HandlerOptions{
			Level: level,
		}))
	}

	return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: level,
	}))
}
