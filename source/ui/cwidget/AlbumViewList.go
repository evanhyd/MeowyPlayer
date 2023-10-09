package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
)

type AlbumViewList = viewList[resource.Album]

func NewAlbumViewList(data *cbinding.AlbumDataList, makeView func(resource.Album) fyne.CanvasObject, size fyne.Size) *AlbumViewList {
	list := &AlbumViewList{display: container.NewGridWrap(size), makeView: makeView}
	data.Attach(list)
	list.ExtendBaseWidget(list)
	return list
}
