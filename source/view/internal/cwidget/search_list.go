package cwidget

import (
	"fmt"
	"playground/pattern"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WidgetObserver[T any] interface {
	fyne.Widget
	pattern.Observer[T]
}

type SearchList[DataType any, WidgetType WidgetObserver[DataType]] struct {
	widget.BaseWidget

	//top
	dropDown    *DropDown
	searchEntry *widget.Entry
	toolbar     *widget.Toolbar

	//main
	scroll *container.Scroll

	itemContainer   *fyne.Container
	itemConstructor func() WidgetType
	data            []DataType
}

func NewSearchList[DataType any, WidgetType WidgetObserver[DataType]](
	itemContainer *fyne.Container,
	itemConstructor func() WidgetType,
	onTextChanged func(string),
	onTextSubmitted func(string),
) *SearchList[DataType, WidgetType] {
	l := SearchList[DataType, WidgetType]{
		dropDown:        NewDropDown(),
		searchEntry:     widget.NewEntry(),
		toolbar:         widget.NewToolbar(),
		scroll:          container.NewScroll(itemContainer),
		itemContainer:   itemContainer,
		itemConstructor: itemConstructor,
	}
	l.searchEntry.OnChanged = onTextChanged
	l.searchEntry.OnSubmitted = onTextSubmitted

	l.ExtendBaseWidget(&l)
	return &l
}

func (v *SearchList[DataType, WidgetType]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewBorder(nil, nil, v.dropDown, v.toolbar, v.searchEntry),
		nil,
		nil,
		nil,
		v.scroll,
	))
}

func (v *SearchList[DataType, WidgetType]) AddDropDown(item *fyne.MenuItem) {
	v.dropDown.Add(item)
}

func (v *SearchList[DataType, WidgetType]) AddToolbar(item fyne.Widget) {
	v.toolbar.Append(&toolbarWidget{item})
}

func (v *SearchList[DataType, WidgetType]) ClearSearchEntry() {
	v.searchEntry.Text = ""
	if v.searchEntry.OnChanged != nil {
		v.searchEntry.OnChanged("")
	}
	v.searchEntry.Refresh()
}

func (v *SearchList[DataType, WidgetType]) Update(data []DataType) {
	fmt.Println("update")
	v.data = data

	//resize to fit, keep the capacity
	if need := len(data) - len(v.itemContainer.Objects); need <= 0 {
		clear(v.itemContainer.Objects[len(data):])
		v.itemContainer.Objects = v.itemContainer.Objects[:len(data)]
	} else {
		for i := 0; i < need; i++ {
			v.itemContainer.Objects = append(v.itemContainer.Objects, v.itemConstructor())
		}
	}

	//update content
	for i := range data {
		v.itemContainer.Objects[i].(WidgetObserver[DataType]).Notify(data[i])
	}

	//update layout
	if v.itemContainer.Layout != nil {
		v.itemContainer.Layout.Layout(v.itemContainer.Objects, v.itemContainer.Size())
	}

	v.scroll.ScrollToTop()
}

var cnt = 0

func (v *SearchList[DataType, WidgetType]) Refresh() {
	v.BaseWidget.Refresh()
	cnt++
	fmt.Println("refresh", cnt)
}
