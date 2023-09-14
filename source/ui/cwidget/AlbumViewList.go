package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type AlbumViewList struct {
	widget.BaseWidget
	grid     *fyne.Container
	makeView func(*player.Album) fyne.CanvasObject
}

func NewAlbumViewList(makeAlbumView func(*player.Album) fyne.CanvasObject, size fyne.Size) *AlbumViewList {
	list := &AlbumViewList{widget.BaseWidget{}, container.NewGridWrap(size), makeAlbumView}
	list.ExtendBaseWidget(list)
	return list
}

func (a *AlbumViewList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewScroll(a.grid))
}

func (a *AlbumViewList) Notify(data []player.Album) {
	a.grid.RemoveAll()
	views := a.makeViews(data)
	for i := range views {
		a.grid.Add(views[i])
	}
}

func (a *AlbumViewList) makeViews(data []player.Album) []fyne.CanvasObject {
	views := make([]fyne.CanvasObject, len(data))
	for i := range data {
		views[i] = a.makeView(&data[i])
	}
	return views
}
