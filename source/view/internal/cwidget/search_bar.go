package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SearchBar[T any] struct {
	widget.BaseWidget
	dropDown    *DropDown
	searchEntry *widget.Entry
	toolbar     *widget.Toolbar
	content     WidgetObserver[T]
}

func NewSearchBar[T any](content WidgetObserver[T], onTextChanged func(string), onTextSubmitted func(string)) *SearchBar[T] {
	l := SearchBar[T]{dropDown: NewDropDown(), searchEntry: widget.NewEntry(), toolbar: widget.NewToolbar(), content: content}
	l.searchEntry.OnChanged = onTextChanged
	l.searchEntry.OnSubmitted = onTextSubmitted
	l.ExtendBaseWidget(&l)
	return &l
}

func (v *SearchBar[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewBorder(nil, nil, v.dropDown, v.toolbar, v.searchEntry), nil,
		nil, nil,
		v.content,
	))
}

func (v *SearchBar[T]) AddDropDown(item *fyne.MenuItem) {
	v.dropDown.Add(item)
}

func (v *SearchBar[T]) AddToolbar(item fyne.Widget) {
	v.toolbar.Append(&toolbarWidget{item})
}

func (v *SearchBar[T]) ClearSearchEntry() {
	v.searchEntry.Text = ""
	if v.searchEntry.OnChanged != nil {
		v.searchEntry.OnChanged("")
	}
	v.searchEntry.Refresh()
}

func (v *SearchBar[T]) Update(data T) {
	v.content.Notify(data)
}

func (v *SearchBar[T]) Refresh() {
	v.BaseWidget.Refresh()
}
