package media

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
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
	imgSource io.Reader,
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
