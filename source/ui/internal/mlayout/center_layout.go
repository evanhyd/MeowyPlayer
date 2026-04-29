package mlayout

import (
	"fyne.io/fyne/v2"
	"golang.org/x/exp/slog"
)

var _ fyne.Layout = &CenterLayout{}

// Fixate the object in the middle of the container, and expand horizontal and vertically by ratio.
type CenterLayout struct {
	hRatio float32
	vRatio float32
}

func NewCenterLayout(horizontalRatio float32, verticalRatio float32) *CenterLayout {
	return &CenterLayout{horizontalRatio, verticalRatio}
}

func (l *CenterLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 1 {
		slog.Error("center layout must have 1 object")
	}

	w := containerSize.Width * l.hRatio
	h := containerSize.Height * l.vRatio
	wOffset := (containerSize.Width - w) / 2
	hOffset := (containerSize.Height - h) / 2
	objects[0].Resize(fyne.NewSize(w, h))
	objects[0].Move(fyne.NewPos(wOffset, hOffset))
}

func (l *CenterLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return objects[0].MinSize()
}
