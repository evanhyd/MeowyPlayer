package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/utility/pattern"
)

type ViewList[T any] struct {
	widget.BaseWidget
	items    *fyne.Container
	scroll   *container.Scroll
	makeView func(T) fyne.CanvasObject
}

func NewViewList[T any](dataList pattern.Subject[[]T], items *fyne.Container, makeView func(T) fyne.CanvasObject) *ViewList[T] {
	viewList := &ViewList[T]{items: items, scroll: container.NewScroll(items), makeView: makeView}
	viewList.ExtendBaseWidget(viewList)
	dataList.Attach(viewList)
	return viewList
}

func (v *ViewList[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.scroll)
}

func (v *ViewList[T]) Notify(data []T) {
	v.items.RemoveAll()
	for i := range data {
		v.items.Add(v.makeView(data[i]))
	}
	v.scroll.Offset = fyne.NewPos(0, 0)
	v.Refresh()
}
