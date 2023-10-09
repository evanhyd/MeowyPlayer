package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
)

type MusicViewList = ViewList[resource.Music]

func NewMusicViewList(data *cbinding.MusicDataList, makeView func(resource.Music) fyne.CanvasObject) *MusicViewList {
	list := &MusicViewList{display: container.NewVBox(), makeView: makeView}
	data.Attach(list)
	list.ExtendBaseWidget(list)
	return list
}
