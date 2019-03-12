package imageutils

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"

	// Imported for its initialization side-effects
	_ "image/png"

	"github.com/nfnt/resize"
)

// OpenImage opens and decode image file
func OpenImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return img, err
}

// WriteJpegImage encodes image as JPEG and write it to the file
func WriteJpegImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, nil)
}

// ResizeImage fits decoded image to maxSize by a longest side
func ResizeImage(maxSize uint, decoded image.Image) image.Image {
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

// OverlayImage tiles a centered overlay image over the base
func OverlayImage(base image.Image, overlay image.Image) image.Image {
	baseBounds := base.Bounds()

	overlayed := image.NewRGBA(baseBounds)
	draw.Draw(overlayed, baseBounds, base, image.ZP, draw.Src)
	draw.Draw(overlayed, baseBounds, overlay, image.ZP, draw.Over)

	return overlayed
}

// ResizeNWatermark resizes a base image to the maxSize,
// then add a watermark fill from the tile
func ResizeNWatermark(base, tile image.Image, maxSize uint) (image.Image, error) {
	resized := ResizeImage(maxSize, base)
	watermark, err := newTiledImage(tile, resized.Bounds(), cc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return OverlayImage(resized, watermark), nil
}

// ceilToOdd rounds x up to nearest odd number
func ceilToOdd(x float64) int {
	ceil := int(math.Ceil(x))

	if isOdd(ceil) {
		return ceil
	}

	return ceil + 1
}

// isOdd test x is odd number
func isOdd(x int) bool {
	return x%2 != 0
}

// tileNum returns a number of tile images in the needed size
func tileNum(boundSize int, tileSize int) int {
	return ceilToOdd(float64(boundSize) / float64(tileSize))
}

// getInitCoord returns start co-ordinate for image tiling in bounds
func getInitCoord(boundWidth int, tileWidth int, n int) int {
	return (boundWidth - n*tileWidth) / 2
}

// align is an image aligns enum
type align int

const (
	tl align = iota
	tc
	tr
	cl
	cc
	cr
	bl
	bc
	br
)

// newTiledImage returns a new image tiled with img in b bounds
func newTiledImage(img image.Image, b image.Rectangle, a align) (image.Image, error) {
	tiled := image.NewRGBA(b)
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	if a != cc {
		return nil, errors.New("Alignments except align.cc are not implemented")
	}

	tiled = image.NewRGBA(b)

	// Number of images on the axes
	xNum := tileNum(b.Dx(), w)
	yNum := tileNum(b.Dy(), h)

	// Initial co-ordinates
	x := getInitCoord(b.Dx(), w, xNum)
	y := getInitCoord(b.Dy(), h, yNum)

	for kx := 0; kx < xNum; kx++ {
		for ky := 0; ky < yNum; ky++ {
			p := image.Point{x + kx*w, y + ky*h}
			draw.Draw(tiled, img.Bounds(), img, p, draw.Over)
		}
	}

	return tiled, nil
}
