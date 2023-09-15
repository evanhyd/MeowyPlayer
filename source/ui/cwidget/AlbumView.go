package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type AlbumView struct {
	widget.BaseWidget
	cover             *canvas.Image
	info              *widget.Label
	title             *widget.Label
	name              string
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func NewAlbumView(album *player.Album) *AlbumView {
	view := &AlbumView{
		widget.BaseWidget{},
		canvas.NewImageFromResource(album.Cover),
		widget.NewLabel(album.Description()),
		widget.NewLabel(album.Title),
		album.Title,
		func(*fyne.PointEvent) {},
		func(*fyne.PointEvent) {},
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
}

func (a *AlbumView) Tapped(event *fyne.PointEvent) {
	a.OnTapped(event)
}

func (a *AlbumView) TappedSecondary(event *fyne.PointEvent) {
	a.OnTappedSecondary(event)
}
