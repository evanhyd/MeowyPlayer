package cwidget

import (
	"playground/pattern"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ObserverCanvasObject[T any] interface {
	fyne.CanvasObject
	pattern.Observer[T]
}

type ScrollList[T any] struct {
	widget.BaseWidget
	scroll      *container.Scroll
	structure   *fyne.Container
	constructor func() ObserverCanvasObject[T]
}

func NewScrollList[T any](structure *fyne.Container, constructor func() ObserverCanvasObject[T]) *ScrollList[T] {
	v := &ScrollList[T]{scroll: container.NewScroll(structure), structure: structure, constructor: constructor}
	v.ExtendBaseWidget(v)
	return v
}

func (v *ScrollList[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.scroll)
}

func (v *ScrollList[T]) Notify(data []T) {
	//resize to fit, keep the capacity
	if len(data) < len(v.structure.Objects) {
		clear(v.structure.Objects[len(data):])
		v.structure.Objects = v.structure.Objects[:len(data)]
	} else if len(data) > len(v.structure.Objects) {
		required := len(data) - len(v.structure.Objects)
		for i := 0; i < required; i++ {
			v.structure.Objects = append(v.structure.Objects, v.constructor())
		}
	}

	//update content
	for i := range data {
		v.structure.Objects[i].(ObserverCanvasObject[T]).Notify(data[i])
	}

	//update layout
	if v.structure.Layout != nil {
		v.structure.Layout.Layout(v.structure.Objects, v.structure.Size())
	}

	v.scroll.ScrollToTop()
}
