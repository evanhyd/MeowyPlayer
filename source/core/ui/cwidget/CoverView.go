package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/resource"
)

type CoverView struct {
	widget.BaseWidget
	tappableBase
	cover *canvas.Image
	title *widget.Label
}

func NewCoverView(size fyne.Size) *CoverView {
	view := &CoverView{
		cover: canvas.NewImageFromResource(resource.DefaultIcon),
		title: widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
	}
	view.cover.SetMinSize(size)
	view.title.Wrapping = fyne.TextWrapWord
	view.title.Hide()
	view.ExtendBaseWidget(view)
	return view
}

func (v *CoverView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(v.cover, v.title))
}

func (v *CoverView) SetAlbum(album *resource.Album) {
	v.cover.Resource = album.Cover
	v.title.SetText(album.Title)
	v.Refresh()
}

func (v *CoverView) MouseIn(*desktop.MouseEvent) {
	v.cover.Translucency = 0.8
	v.title.Show()
}

func (v *CoverView) MouseOut() {
	v.cover.Translucency = 0.0
	v.title.Hide()
}

func (v *CoverView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
