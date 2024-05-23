package cwidget

import (
	"image/color"
	"playground/pattern"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type itemHitbox struct {
	widget.BaseWidget
	TappableBase
}

func newHitbox() *itemHitbox {
	h := &itemHitbox{}
	h.ExtendBaseWidget(h)
	return h
}

func (i *itemHitbox) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

func (i *itemHitbox) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type WidgetObserver[T any] interface {
	fyne.Widget
	pattern.Observer[T]
}

type item[T any] struct {
	widget.BaseWidget
	*itemHitbox
	content WidgetObserver[T]
}

func newItem[T any](content WidgetObserver[T]) *item[T] {
	i := &item[T]{itemHitbox: newHitbox(), content: content}
	i.ExtendBaseWidget(i)
	return i
}

func (i *item[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(i.content, i.itemHitbox))
}

func (i *item[T]) Notify(data T) {
	i.content.Notify(data)
}

type ItemTapEvent[T any] struct {
	*fyne.PointEvent
	Data T
}

type ScrollList[T any] struct {
	widget.BaseWidget
	scroll                *container.Scroll
	display               *fyne.Container
	makeItem              func() WidgetObserver[T]
	OnItemTapped          func(ItemTapEvent[T])
	OnItemTappedSecondary func(ItemTapEvent[T])
}

func NewScrollList[T any](display *fyne.Container, makeItem func() WidgetObserver[T]) *ScrollList[T] {
	v := &ScrollList[T]{scroll: container.NewScroll(display), display: display, makeItem: makeItem}
	v.ExtendBaseWidget(v)
	return v
}

func (v *ScrollList[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.scroll)
}

func (v *ScrollList[T]) Notify(data []T) {

	//resize to fit, keep the capacity
	if need := len(data) - len(v.display.Objects); need <= 0 {
		clear(v.display.Objects[len(data):])
		v.display.Objects = v.display.Objects[:len(data)]
	} else {
		for i := 0; i < need; i++ {
			i := len(v.display.Objects)
			item := newItem(v.makeItem())
			item.OnTapped = func(e *fyne.PointEvent) { v.OnItemTapped(ItemTapEvent[T]{e, data[i]}) }
			item.OnTappedSecondary = func(e *fyne.PointEvent) { v.OnItemTappedSecondary(ItemTapEvent[T]{e, data[i]}) }
			v.display.Objects = append(v.display.Objects, item)
		}
	}

	//update content
	for i := range data {
		v.display.Objects[i].(*item[T]).Notify(data[i])
	}

	//update layout
	if v.display.Layout != nil {
		v.display.Layout.Layout(v.display.Objects, v.display.Size())
	}

	v.scroll.ScrollToTop()
}
