package imageutils

import (
	"image"

	"github.com/nfnt/resize"
)

func resizeImage(maxSize uint, decoded image.Image) image.Image {
	width := decoded.Bounds().Max.X
	height := decoded.Bounds().Max.Y

	var targetWidth uint
	var targetHeight uint

	if width > height {
		targetWidth = maxSize
	} else {
		targetHeight = maxSize
	}

	return resize.Resize(targetWidth, targetHeight, decoded, resize.Lanczos3)
}
