package mlayout

import (
	"fyne.io/fyne/v2"
	"golang.org/x/exp/slog"
)

// Split the canvas into up and down determine by the ratio.
type VSplitLayout struct {
	ratio float32
}

func NewVSplitLayout(ratio float32) *VSplitLayout {
	return &VSplitLayout{ratio: ratio}
}

func (l *VSplitLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 2 {
		slog.Error("h-split layout must have 2 objects")
	}

	h1 := containerSize.Height * l.ratio
	h2 := containerSize.Height - h1
	w := containerSize.Width
	objects[0].Resize(fyne.NewSize(w, h1))
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Resize(fyne.NewSize(w, h2))
	objects[1].Move(fyne.NewPos(0, h1))
}

func (l *VSplitLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(max(objects[0].MinSize().Width, objects[1].MinSize().Width), objects[0].MinSize().Height+objects[1].MinSize().Height)
}
