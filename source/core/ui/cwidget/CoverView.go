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

func (c *CoverView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewMax(c.cover, c.title))
}

func (c *CoverView) SetAlbum(album *resource.Album) {
	c.cover.Resource = album.Cover
	c.title.SetText(album.Title)
	c.Refresh()
}

func (c *CoverView) MouseIn(*desktop.MouseEvent) {
	c.cover.Translucency = 0.8
	c.title.Show()
}

func (c *CoverView) MouseOut() {
	c.cover.Translucency = 0.0
	c.title.Hide()
}

func (c *CoverView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
