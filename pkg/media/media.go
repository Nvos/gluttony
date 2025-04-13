package media

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"os"
)

type Store struct {
	rootDir *os.Root
}

func NewStore(rootDir *os.Root) *Store {
	return &Store{rootDir: rootDir}
}

func (s *Store) UploadImage(file io.Reader) (string, error) {
	fileName := fmt.Sprintf("%s.webp", uuid.New().String())
	create, err := s.rootDir.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("upload image file: %w", err)
	}
	defer func(create fs.File) {
		_ = create.Close()
	}(create)

	err = optimizeAndWriteImage(file, create)
	if err != nil {
		return "", fmt.Errorf("upload iamge (optimizing): %w", err)
	}

	if err := create.Sync(); err != nil {
		return "", fmt.Errorf("upload image (file sync): %w", err)
	}

	return fileName, nil
}
