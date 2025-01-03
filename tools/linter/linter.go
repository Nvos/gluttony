package linter

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const downloadURL = "https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.zip"

type Linter struct {
	version string
	binPath string
	logger  *slog.Logger
}

func NewLinter(version string, binPath string, logger *slog.Logger) *Linter {
	return &Linter{version: version, binPath: binPath, logger: logger}
}

func (l *Linter) Run(ctx context.Context) error {
	if err := l.installGolangCiLint(ctx); err != nil {
		return fmt.Errorf("install golangci-lint: %w", err)
	}

	execPath := filepath.Join(l.binPath, "golangci-lint.exe")
	cmd := exec.CommandContext(ctx, execPath, "run", "./...")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (l *Linter) installGolangCiLint(ctx context.Context) error {
	filename := "golangci-lint"
	if runtime.GOOS == "windows" {
		filename = "golangci-lint.exe"
	}

	outFilePath := filepath.Join(l.binPath, filename)
	downloadURL := fmt.Sprintf(
		downloadURL, l.version,
		l.version,
		runtime.GOOS,
		runtime.GOARCH,
	)

	if _, err := os.Stat(outFilePath); !os.IsNotExist(err) {
		return nil
	}

	l.logger.InfoContext(ctx, "downloading golangci-lint", slog.String("version", l.version))
	create, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer create.Close()

	resp, err := download(ctx, downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	for i := range zipReader.File {
		file := zipReader.File[i]
		if strings.HasPrefix(filepath.Base(file.Name), "golangci-lint") {
			if err := writeZipFile(create, file); err != nil {
				return err
			}
		}
	}

	return nil
}
