package widgets

import (
	"meowyplayer/storages"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = &PlaylistCard{}

type PlaylistCard struct {
	widget.BaseWidget
	cover     *canvas.Image
	title     *widget.Label
	highlight *canvas.Rectangle
	playlist  storages.Playlist
	onTapped  func(playlist storages.Playlist)
}

func NewPlaylistCard(onTapped func(playlist storages.Playlist)) *PlaylistCard {
	p := PlaylistCard{
		cover:     canvas.NewImageFromResource(nil),
		title:     widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		highlight: canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
		onTapped:  onTapped,
	}
	p.cover.SetMinSize(PlaylistCardSize)
	p.title.Wrapping = fyne.TextWrapWord
	p.highlight.Hide()
	p.ExtendBaseWidget(&p)
	return &p
}

func (p *PlaylistCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, p.title, nil, nil, p.cover),
		p.highlight,
	))
}

func (p *PlaylistCard) MouseIn(*desktop.MouseEvent) {
	p.highlight.Show()
	p.Refresh()
}

func (p *PlaylistCard) MouseOut() {
	p.highlight.Hide()
	p.Refresh()
}

func (p *PlaylistCard) MouseMoved(*desktop.MouseEvent) {
	// Satisfy desktop.Hoverable
}

func (p *PlaylistCard) Tapped(*fyne.PointEvent) {
	p.onTapped(p.playlist)
}

func (p *PlaylistCard) Set(playlist storages.Playlist) {
	p.cover.Resource = fyne.NewStaticResource(strconv.FormatInt(playlist.PlaylistId, 16), playlist.CoverBlob)
	p.title.SetText(playlist.Title)
	p.playlist = playlist
	p.Refresh()
}
