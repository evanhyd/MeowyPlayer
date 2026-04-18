package widgets

import (
	"bytes"
	"image"
	"io"

	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2"
	"golang.org/x/image/draw"
)

func ScaleImageFromReader(reader io.Reader, size fyne.Size) (image.Image, error) {
	originalThumbnail, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Create a new RGBA image of the target size
	scaledThumbnail := image.NewRGBA(image.Rect(0, 0, int(size.Width), int(size.Height)))

	// Use draw.CatmullRom for high-quality scaling
	draw.CatmullRom.Scale(
		scaledThumbnail, scaledThumbnail.Rect,
		originalThumbnail, originalThumbnail.Bounds(),
		draw.Over, nil,
	)
	return scaledThumbnail, nil
}

func ScaleImageFromBytes(data []byte, size fyne.Size) (image.Image, error) {
	return ScaleImageFromReader(bytes.NewBuffer(data), size)
}
