package view

import (
	"fmt"
	"playground/browser"
	"playground/model"
	"playground/resource"
	"playground/view/internal/cwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type HomePage struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[model.Music]
	browser   browser.Browser
}

func NewHomePage(client *model.Client) *HomePage {
	var v HomePage
	v = HomePage{
		searchBar: cwidget.NewSearchBar[model.Music](v.render),
		browser:   browser.NewYouTubeBrowser(),
	}
	v.searchBar.AddMenuItem("YouTube", resource.YouTubeIcon, func() {})
	v.searchBar.Select(0)

	v.ExtendBaseWidget(&v)
	return &v
}

func (v *HomePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(v.searchBar, nil, nil, nil))
}

func (v *HomePage) render() {
	fmt.Println("render")
}
