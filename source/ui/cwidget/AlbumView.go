package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/resource"
)

type AlbumView struct {
	widget.BaseWidget
	tappableBase
	cover *canvas.Image
	info  *widget.Label
	title *widget.Label
	name  string
}

func NewAlbumView(album *resource.Album) *AlbumView {
	view := &AlbumView{
		cover: canvas.NewImageFromResource(album.Cover),
		info:  widget.NewLabel(album.Description()),
		title: widget.NewLabel(album.Title),
		name:  album.Title,
	}
	view.info.Hide()
	view.info.Wrapping = fyne.TextWrapWord
	view.title.Alignment = fyne.TextAlignCenter
	view.title.Wrapping = fyne.TextTruncate
	view.ExtendBaseWidget(view)
	return view
}

func (a *AlbumView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, a.title, nil, nil, container.NewMax(a.cover, a.info)))
}

func (a *AlbumView) MouseIn(event *desktop.MouseEvent) {
	a.title.Text = ""
	a.cover.Translucency = 0.8
	a.info.Show()
	a.Refresh()
}

func (a *AlbumView) MouseOut() {
	a.title.Text = a.name
	a.cover.Translucency = 0.0
	a.info.Hide()
	a.Refresh()
}

func (a *AlbumView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
