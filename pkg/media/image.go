package media

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
	"mime/multipart"
	"net/http"
	// Register webp codec.
	_ "golang.org/x/image/webp"
	"image"
	// Register jpeg codec.
	_ "image/jpeg"
	// Register png codec.
	_ "image/png"
	"io"
)

func optimizeAndWriteImage(
	imgSource multipart.File,
	writeTo io.Writer,
) error {
	img, _, err := image.Decode(imgSource)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	const expectedWidth = 382
	const expectedHeight = 256
	img = imaging.Fill(img, expectedWidth, expectedHeight, imaging.Center, imaging.Lanczos)

	const defaultImageQuality = 80
	const defaultMethod = 4
	encodeOpts := webp.Options{
		Quality:  defaultImageQuality,
		Method:   defaultMethod,
		Lossless: false,
		Exact:    false,
	}
	if err := webp.Encode(writeTo, img, encodeOpts); err != nil {
		return fmt.Errorf("encode image to webp: %w", err)
	}

	return nil
}

func isMediaContentAllowed(file multipart.File) (bool, error) {
	buffer := make([]byte, 512)

	if _, err := file.Read(buffer); err != nil {
		return false, fmt.Errorf("read header bytes: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return false, fmt.Errorf("seek to start: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return true, nil
	default:
		return false, nil
	}

}
