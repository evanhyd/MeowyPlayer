package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/resource"
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
		title: widget.NewLabelWithStyle(album.Title, fyne.TextAlignCenter, fyne.TextStyle{}),
		name:  album.Title,
	}
	view.info.Hide()
	view.info.Wrapping = fyne.TextWrapWord
	view.title.Wrapping = fyne.TextTruncate
	view.ExtendBaseWidget(view)
	return view
}

func (v *AlbumView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, v.title, nil, nil, container.NewStack(v.cover, v.info)))
}

func (v *AlbumView) MouseIn(event *desktop.MouseEvent) {
	v.title.SetText("")
	v.cover.Translucency = 0.8
	v.info.Show()
}

func (v *AlbumView) MouseOut() {
	v.title.SetText(v.name)
	v.cover.Translucency = 0.0
	v.info.Hide()
}

func (v *AlbumView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
