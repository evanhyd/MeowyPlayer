package cwidget

import (
	"fyne.io/fyne/v2"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/utility/network/scraper"
)

type YoutubeResultView struct {
}

func NewYoutubeResultView(data *cbinding.YoutubeResultDataList, makeView func(*scraper.YoutubeResult) fyne.CanvasObject) *YoutubeResultView {
	return nil
}
