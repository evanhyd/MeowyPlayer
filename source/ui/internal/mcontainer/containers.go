package mcontainer

import (
	"meowyplayer/ui/internal/mlayout"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewCenter(horizontalRatio float32, verticalRatio float32, object fyne.CanvasObject) *fyne.Container {
	return container.New(mlayout.NewCenterLayout(horizontalRatio, verticalRatio), object)
}

func NewHSplit(ratio float32, left fyne.CanvasObject, right fyne.CanvasObject) *fyne.Container {
	return container.New(mlayout.NewHSplitLayout(ratio), left, right)
}

func NewVSplit(ratio float32, up fyne.CanvasObject, down fyne.CanvasObject) *fyne.Container {
	return container.New(mlayout.NewVSplitLayout(ratio), up, down)
}
