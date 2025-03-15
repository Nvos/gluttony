package media

import (
	"fmt"
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

type optimizeImageOpts struct {
	quality int
}

func optimizeAndWriteImage(
	imgSource io.Reader,
	writeTo io.Writer,
	opts optimizeImageOpts,
) error {
	img, _, err := image.Decode(imgSource)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	// TODO: crop + resize

	encodeOpts := webp.Options{
		Quality:  opts.quality,
		Lossless: false,
		Method:   0,
		Exact:    false,
	}
	if err := webp.Encode(writeTo, img, encodeOpts); err != nil {
		return fmt.Errorf("encode image to webp: %w", err)
	}

	return nil
}
