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
}

func NewPlaylistCard(onTapped func(playlist storages.Playlist)) *PlaylistCard {
	c := PlaylistCard{
		cover:     canvas.NewImageFromResource(nil),
		title:     widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		highlight: canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
	}
	c.cover.SetMinSize(PlaylistCardSize)
	c.title.Wrapping = fyne.TextWrapWord
	c.highlight.Hide()
	c.ExtendBaseWidget(&c)
	return &c
}

func (c *PlaylistCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, c.title, nil, nil, c.cover),
		c.highlight,
	))
}

func (c *PlaylistCard) MouseIn(*desktop.MouseEvent) {
	c.highlight.Show()
	c.Refresh()
}

func (c *PlaylistCard) MouseOut() {
	c.highlight.Hide()
	c.Refresh()
}

func (c *PlaylistCard) MouseMoved(*desktop.MouseEvent) {
	// Satisfy desktop.Hoverable
}

func (b *PlaylistCard) Tapped(*fyne.PointEvent) {
	// Disable the yellow highlight from widget.List.
}

func (c *PlaylistCard) Set(playlist storages.Playlist) {
	c.cover.Resource = fyne.NewStaticResource(strconv.FormatInt(playlist.PlaylistId, 16), playlist.CoverBlob)
	c.title.SetText(playlist.Title)
	c.playlist = playlist
}
