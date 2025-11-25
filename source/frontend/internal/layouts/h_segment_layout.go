package layouts

import "fyne.io/fyne/v2"

type HSegmentLayout struct {
	widthRatio []float32
}

func NewHSegmentLayout(widthRatio ...float32) *HSegmentLayout {
	totalWidth := float32(0)
	for _, width := range widthRatio {
		totalWidth += width
	}
	for i, width := range widthRatio {
		widthRatio[i] = width / totalWidth
	}
	return &HSegmentLayout{widthRatio}
}

func (l *HSegmentLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != len(l.widthRatio) {
		panic("object lengths doesn't match up with the widthRatio length")
	}

	pos := fyne.NewPos(0, 0)
	for i, object := range objects {
		width := containerSize.Width * l.widthRatio[i]
		height := object.MinSize().Height
		object.Resize(fyne.NewSize(width, height))
		object.Move(pos)
		pos = pos.AddXY(width, 0)
	}
}

func (l *HSegmentLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, o := range objects {
		childSize := o.MinSize()
		minSize = minSize.AddWidthHeight(childSize.Width, childSize.Height)
	}
	return minSize
}
