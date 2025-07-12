package image

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
)

const maxFileSize = 5242880 // 5MB

type Service struct {
	rootDir *os.Root
}

func NewService(rootDir *os.Root) *Service {
	return &Service{rootDir: rootDir}
}

func (s *Service) Delete(fileName string) error {
	if err := s.rootDir.Remove(fileName); err != nil {
		return fmt.Errorf("delete image by filename='%s': %w", fileName, err)
	}

	return nil
}

func (s *Service) Upload(file *multipart.FileHeader) (string, error) {
	if file.Size > maxFileSize {
		return "", errors.New("file size is too big")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open image file: %w", err)
	}
	defer func(src io.ReadCloser) {
		_ = src.Close()
	}(src)

	ok, err := isMediaContentAllowed(src)
	if err != nil {
		return "", fmt.Errorf("is media content allowed: %w", err)
	}
	if !ok {
		return "", errors.New("media content is not allowed")
	}

	fileName := fmt.Sprintf("%s.webp", uuid.New().String())
	create, err := s.rootDir.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("upload image file: %w", err)
	}
	defer func(create fs.File) {
		_ = create.Close()
	}(create)

	err = optimizeAndWriteImage(src, create)
	if err != nil {
		return "", fmt.Errorf("upload iamge (optimizing): %w", err)
	}

	return fileName, nil
}
