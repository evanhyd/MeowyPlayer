package cwidget

import (
	"playground/pattern"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WidgetObserver[T any] interface {
	fyne.Widget
	pattern.Observer[T]
}

type cachedList[T any] struct {
	widget.List
	data []T
}

func NewCachedList[T any, WidgetType WidgetObserver[T]](itemConstructor func() WidgetType) *cachedList[T] {
	var l cachedList[T]
	l = cachedList[T]{List: widget.List{
		Length:     func() int { return len(l.data) },
		CreateItem: func() fyne.CanvasObject { return itemConstructor() },
		UpdateItem: func(i widget.ListItemID, item fyne.CanvasObject) { item.(WidgetObserver[T]).Notify(l.data[i]) },
	}}
	l.ExtendBaseWidget(&l)
	return &l
}

func (l *cachedList[T]) Notify(data []T) {
	l.data = data
	l.Refresh()
}

type customList[T any] struct {
	container.Scroll
	itemContainer   *fyne.Container
	itemConstructor func() fyne.Widget
	data            []T
}

func NewCustomList[T any, WidgetType WidgetObserver[T]](itemContainer *fyne.Container, itemConstructor func() WidgetType) *customList[T] {
	var l customList[T]
	l = customList[T]{
		Scroll:          container.Scroll{Direction: container.ScrollBoth, Content: itemContainer},
		itemContainer:   itemContainer,
		itemConstructor: func() fyne.Widget { return itemConstructor() },
	}
	l.ExtendBaseWidget(&l)
	return &l
}

func (l *customList[T]) Notify(data []T) {
	l.data = data

	//resize to fit, keep the capacity
	if need := len(data) - len(l.itemContainer.Objects); need <= 0 {
		clear(l.itemContainer.Objects[len(data):])
		l.itemContainer.Objects = l.itemContainer.Objects[:len(data)]
	} else {
		for i := 0; i < need; i++ {
			l.itemContainer.Objects = append(l.itemContainer.Objects, l.itemConstructor())
		}
	}

	//update content
	for i := range data {
		l.itemContainer.Objects[i].(WidgetObserver[T]).Notify(data[i])
	}

	//update layout
	if l.itemContainer.Layout != nil {
		l.itemContainer.Layout.Layout(l.itemContainer.Objects, l.itemContainer.Size())
	}

	l.ScrollToTop()
}
