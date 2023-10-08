package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type viewList[T any] struct {
	widget.BaseWidget
	display  *fyne.Container
	makeView func(*T) fyne.CanvasObject
}

func (v *viewList[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewScroll(v.display))
}

func (v *viewList[T]) Notify(data []T) {
	v.display.RemoveAll()
	views := v.makeViews(data)
	for _, view := range views {
		v.display.Add(view)
	}
}

func (v *viewList[T]) makeViews(data []T) []fyne.CanvasObject {
	views := make([]fyne.CanvasObject, len(data))
	for i := range data {
		views[i] = v.makeView(&data[i])
	}
	return views
}
