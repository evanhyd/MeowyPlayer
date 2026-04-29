package mlayout

import (
	"fyne.io/fyne/v2"
	"golang.org/x/exp/slog"
)

// Split the canvas into left and right determine by the ratio.
type HSplitLayout struct {
	ratio float32
}

func NewHSplitLayout(ratio float32) *HSplitLayout {
	return &HSplitLayout{ratio: ratio}
}

func (l *HSplitLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 2 {
		slog.Error("h-split layout must have 2 objects")
	}

	w1 := containerSize.Width * l.ratio
	w2 := containerSize.Width - w1
	h := containerSize.Height
	objects[0].Resize(fyne.NewSize(w1, h))
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Resize(fyne.NewSize(w2, h))
	objects[1].Move(fyne.NewPos(w1, 0))
}

func (l *HSplitLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(objects[0].MinSize().Width+objects[1].MinSize().Width,
		max(objects[0].MinSize().Height, objects[1].MinSize().Height))
}
