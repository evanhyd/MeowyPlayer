package mwidget

import (
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mutil"

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
	highlight         *canvas.Rectangle
	cover             *canvas.Image
	title             *widget.Label
	playlist          storages.Playlist
	onTapped          func(playlist storages.Playlist)
	onTappedSecondary func(playlist storages.Playlist, event *fyne.PointEvent)
}

func NewPlaylistCard(
	onTapped func(playlist storages.Playlist),
	onTappedSecondary func(playlist storages.Playlist, event *fyne.PointEvent),
) *PlaylistCard {

	p := PlaylistCard{
		highlight:         canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
		cover:             canvas.NewImageFromResource(nil),
		title:             widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		onTapped:          onTapped,
		onTappedSecondary: onTappedSecondary,
	}
	p.highlight.CornerRadius = 8.0
	p.highlight.Hide()
	p.cover.SetMinSize(PlaylistCardSize)
	p.cover.CornerRadius = 8.0
	p.cover.FillMode = canvas.ImageFillCover
	p.title.Truncation = fyne.TextTruncateEllipsis
	p.ExtendBaseWidget(&p)
	return &p
}

func (c *PlaylistCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, c.title, nil, nil, container.NewStack(c.highlight, c.cover)))
}

func (p *PlaylistCard) MouseIn(*desktop.MouseEvent) {
	p.cover.Translucency = 0.3
	p.highlight.Show()
	p.Refresh()
}

func (p *PlaylistCard) MouseOut() {
	p.cover.Translucency = 0.0
	p.highlight.Hide()
	p.Refresh()
}

func (p *PlaylistCard) MouseMoved(*desktop.MouseEvent) {
	// Satisfy desktop.Hoverable
}

func (p *PlaylistCard) Tapped(*fyne.PointEvent) {
	p.onTapped(p.playlist)
}

func (p *PlaylistCard) TappedSecondary(event *fyne.PointEvent) {
	p.onTappedSecondary(p.playlist, event)
}

func (p *PlaylistCard) Set(playlist storages.Playlist) {
	p.cover.Resource = fyne.NewStaticResource(mutil.PlaylistIdToString(playlist.PlaylistId), playlist.CoverBlob)
	p.title.SetText(playlist.Title)
	p.playlist = playlist
	p.Refresh()
}
