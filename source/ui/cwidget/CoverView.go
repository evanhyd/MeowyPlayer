package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

type CoverView struct {
	widget.BaseWidget
	cover    *canvas.Image
	title    *widget.Label
	size     fyne.Size
	onTapped func(*fyne.PointEvent)
}

func NewCoverView(coverSize fyne.Size) *CoverView {
	view := &CoverView{widget.BaseWidget{}, canvas.NewImageFromResource(resource.DefaultIcon()), widget.NewLabel(""), coverSize, func(*fyne.PointEvent) {}}
	view.cover.SetMinSize(coverSize)
	view.title.Alignment = fyne.TextAlignCenter
	view.title.Wrapping = fyne.TextWrapWord
	view.title.Hide()

	view.ExtendBaseWidget(view)
	return view
}

func (c *CoverView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewMax(c.cover, c.title))
}

func (c *CoverView) SetAlbum(album *player.Album) {
	c.cover.Resource = album.Cover
	c.title.Text = album.Title
	c.Refresh()
}

func (c *CoverView) MouseIn(*desktop.MouseEvent) {
	c.cover.Translucency = 0.8
	c.title.Show()
	c.Refresh()
}

func (c *CoverView) MouseOut() {
	c.cover.Translucency = 0.0
	c.title.Hide()
	c.Refresh()
}

func (c *CoverView) MouseMoved(*desktop.MouseEvent) {
}

func (c *CoverView) Tapped(event *fyne.PointEvent) {
	c.onTapped(event)
}

func (c *CoverView) Notify(play *player.Play) {
	c.SetAlbum(play.Album())
	c.onTapped = func(*fyne.PointEvent) { manager.GetCurrentAlbum().Set(play.Album()) }
}
