package view

import (
	"fmt"
	"playground/browser"
	"playground/model"
	"playground/resource"
	"playground/view/internal/cwidget"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type HomePage struct {
	widget.BaseWidget
	list *cwidget.SearchList[model.Music, *MusicCard]

	client   *model.Client
	pipeline DataPipeline[model.Music]
	browser  browser.Browser
}

func NewHomePage(client *model.Client) *HomePage {
	var v *HomePage
	v = &HomePage{
		list: cwidget.NewSearchList[model.Music, *MusicCard](
			container.NewVBox(),
			newMusicCard,
			func(e cwidget.ItemTapEvent[model.Music]) {
				fmt.Println("left click", e.Data)
			},
			func(e cwidget.ItemTapEvent[model.Music]) {
				fmt.Println("rigth click")
			},
			func(sub string) {
				v.pipeline.filter = func(str string) bool {
					return strings.Contains(strings.ToLower(str), strings.ToLower(sub))
				}
				v.updateList()
			},
			nil,
		),
		client: client,
		pipeline: DataPipeline[model.Music]{
			comparator: func(_, _ model.Music) int { return -1 },
			filter:     func(_ string) bool { return true },
		},
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

func (v *HomePage) updateList() {
	fmt.Println("update home page")
}
