package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type CoverView struct {
	widget.BaseWidget
	display  *fyne.Container
	cover    *canvas.Image
	title    *widget.Label
	OnTapped func(*fyne.PointEvent)
}

func NewCoverView(album *player.Album) *CoverView {
	view := &CoverView{
		widget.BaseWidget{},
		container.NewMax(),
		canvas.NewImageFromResource(album.Cover),
		widget.NewLabel(album.Title),
		func(*fyne.PointEvent) {},
	}
	view.title.Alignment = fyne.TextAlignCenter
	view.title.Wrapping = fyne.TextTruncate
	return view
}

func (c *CoverView) Tapped(event *fyne.PointEvent) {
	c.OnTapped(event)
}

func (c *CoverView) MouseIn(*desktop.MouseEvent) {
	c.cover.Translucency = 0.8
	c.display.Add(c.title)
	c.Refresh()
}

func (c *CoverView) MouseOut() {
	c.cover.Translucency = 0.0
	c.display.Remove(c.title)
	c.Refresh()
}

func (c *CoverView) MouseMoved(*desktop.MouseEvent) {
}
