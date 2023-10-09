package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/utility/network/scraper"
)

type VideoScraperViewList = viewList[scraper.VideoScraper]

func NewScraperViewList(data *cbinding.VideoScraperDataList, makeView func(scraper.VideoScraper) fyne.CanvasObject, size fyne.Size) *VideoScraperViewList {
	list := &VideoScraperViewList{display: container.NewGridWrap(size), makeView: makeView}
	data.Attach(list)
	list.ExtendBaseWidget(list)
	return list
}
