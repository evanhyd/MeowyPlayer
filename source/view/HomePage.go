package view

import (
	"fmt"
	"playground/cwidget"
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type HomePage struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[model.Music]
}

func NewHomePage(client *model.MusicClient) *HomePage {
	var v HomePage
	v = HomePage{
		searchBar: cwidget.NewSearchBar[model.Music](v.render),
	}

	// v.searchBar.AddComparator()

	v.ExtendBaseWidget(&v)
	return &v
}

func (v *HomePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		v.searchBar, nil, nil, nil,
	))
}

func (v *HomePage) render() {
	fmt.Println("render")
}
