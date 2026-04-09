package widgets

import (
	"meowyplayer/storages"

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
	cover      *canvas.Image
	title      *widget.RichText
	highlight  *canvas.Rectangle
	playlistID int64
}

func NewPlaylistCard(onTapped func(playlistId int64)) *PlaylistCard {
	c := PlaylistCard{
		cover:     canvas.NewImageFromResource(nil),
		title:     widget.NewRichTextWithText(""),
		highlight: canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
	}
	c.ExtendBaseWidget(&c)
	return &c
}

func (c *PlaylistCard) CreateRenderer() fyne.WidgetRenderer {
	c.cover.SetMinSize(fyne.NewSize(64, 64))
	c.title.Wrapping = fyne.TextWrapWord
	c.highlight.Hide()
	return widget.NewSimpleRenderer(container.NewBorder(nil, c.title, nil, nil, c.cover))
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

func (c *PlaylistCard) Set(playlist storages.Playlist) {
	scaledCover, err := scaleImage(playlist.CoverBlob, 64, 64)
	if err != nil {
		fyne.LogError("Failed to decode or scale thumbnail", err)
		return
	}
	c.cover.Image = scaledCover
	c.title = widget.NewRichTextWithText(playlist.Title)
	c.playlistID = playlist.PlaylistID
}
