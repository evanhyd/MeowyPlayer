package layouts

import (
	"fyne.io/fyne/v2"
	"golang.org/x/exp/slog"
)

var _ fyne.Layout = &CenterLayout{}

type CenterLayout struct {
	horizontalRatio float32
	verticalRatio   float32
}

// Each widthRatio represents the percentage of the horizontal space that the child canvas object occupies.
func NewCenterLayout(horizontalRatio float32, verticalRatio float32) *CenterLayout {
	return &CenterLayout{horizontalRatio, verticalRatio}
}

func (l *CenterLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 1 {
		slog.Error("object lengths doesn't match up with the widthRatio length")
	}

	width := containerSize.Width * l.horizontalRatio
	height := containerSize.Height * l.verticalRatio
	widthOffset := (containerSize.Width - width) / 2
	heightOffset := (containerSize.Height - height) / 2

	objects[0].Resize(fyne.NewSize(width, height))
	objects[0].Move(fyne.NewPos(widthOffset, heightOffset))
}

func (l *CenterLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return objects[0].MinSize()
}
