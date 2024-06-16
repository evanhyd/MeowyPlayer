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
	list *cwidget.SearchList[browser.Result, *ThumbnailCard]

	client  *model.Client
	browser browser.Browser
}

func NewHomePage(client *model.Client) *HomePage {
	var v *HomePage
	v = &HomePage{
		list: cwidget.NewSearchList(
			container.NewVBox(),
			newThumbnailCard,
			func(e cwidget.ItemTapEvent[browser.Result]) {
				fmt.Println("left click", e.Data)
			},
			func(e cwidget.ItemTapEvent[browser.Result]) {
				fmt.Println("rigth click")
			},
			nil,
			func(title string) {
				result, err := v.browser.Search(title)
				if err != nil {
					fyne.LogError("browser search failed", err)
					return
				}
				v.list.Update(result)
			},
		),
		client: client,
	}

	v.list.AddDropDown(cwidget.NewMenuItemWithIcon("YouTube", resource.YouTubeIcon, func() {
		v.browser = browser.NewYouTubeBrowser()
	}))

	v.ExtendBaseWidget(v)
	return v
}

func (v *HomePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.list)
}
