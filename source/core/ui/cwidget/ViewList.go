package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/utility/pattern"
)

type ViewList[T any] struct {
	widget.BaseWidget
	display  *fyne.Container
	scroll   *container.Scroll
	makeView func(T) fyne.CanvasObject
}

func NewViewList[T any](dataList pattern.Subject[[]T], display *fyne.Container, makeView func(T) fyne.CanvasObject) *ViewList[T] {
	viewList := &ViewList[T]{display: display, scroll: container.NewScroll(display), makeView: makeView}
	dataList.Attach(viewList)
	viewList.ExtendBaseWidget(viewList)
	return viewList
}

func (v *ViewList[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.scroll)
}

func (v *ViewList[T]) Notify(data []T) {
	v.display.RemoveAll()
	views := v.makeViews(data)
	for _, view := range views {
		v.display.Add(view)
	}
	v.scroll.Offset = fyne.NewPos(0, 0)
	v.Refresh()
}

func (v *ViewList[T]) makeViews(data []T) []fyne.CanvasObject {
	views := make([]fyne.CanvasObject, len(data))
	for i := range data {
		views[i] = v.makeView(data[i])
	}
	return views
}
