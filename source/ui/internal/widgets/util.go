package widgets

import (
	"bytes"
	"image"
	"io"

	"golang.org/x/image/draw"
)

func ScaleImageFromReader(reader io.Reader, targetWidth int, targetHeight int) (image.Image, error) {
	originalThumbnail, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Create a new RGBA image of the target size
	scaledThumbnail := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Use draw.CatmullRom for high-quality scaling
	draw.CatmullRom.Scale(
		scaledThumbnail, scaledThumbnail.Rect,
		originalThumbnail, originalThumbnail.Bounds(),
		draw.Over, nil,
	)
	return scaledThumbnail, nil
}

func ScaleImageFromBytes(data []byte, targetWidth int, targetHeight int) (image.Image, error) {
	return ScaleImageFromReader(bytes.NewBuffer(data), targetWidth, targetHeight)
}
