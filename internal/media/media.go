package media

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"io"
)

// TODO: rename to image store
type Store struct {
	fs afero.Fs
}

func NewStore(fs afero.Fs) *Store {
	return &Store{fs: fs}
}

func (s *Store) Store(file io.Reader) (string, error) {
	// TODO: pass extension via arg, ideally would want to just convert to webp but
	// there's no lib offering purego solution. Preferably would want to avoid cgo
	fileName := fmt.Sprintf("%s.jpeg", uuid.New().String())

	create, err := s.fs.Create(fileName)
	if err != nil {
		return "", err
	}
	defer create.Close()

	if _, err := io.Copy(create, file); err != nil {
		return "", fmt.Errorf("store media file (copy): %w", err)
	}

	if err := create.Sync(); err != nil {
		return "", fmt.Errorf("store media file (sync): %w", err)
	}

	return fileName, nil
}

func (s *Store) Get(fileName string) (io.Reader, error) {
	return s.fs.Open(fileName)
}
