package media

import (
	"github.com/gen2brain/webp"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
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
		return err
	}

	// TODO: crop + resize

	return webp.Encode(writeTo, img, webp.Options{
		Quality: opts.quality,
	})
}
