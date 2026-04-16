package widgets

import (
	"fmt"
	"log/slog"
	"meowyplayer/scrapers"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = &ThumbnailCard{}

type ThumbnailCard struct {
	widget.BaseWidget
	thumbnail           *canvas.Image
	heading             *widget.Label
	meta                *widget.Label
	highlight           *canvas.Rectangle
	playInBrowserButton *widget.Button
	addToPlaylistButton *widget.Button
	result              scrapers.Result
}

func NewThumbnailCard(playInBrowserCallback func(scrapers.Result), addToPlaylistCallback func(scrapers.Result)) *ThumbnailCard {
	c := ThumbnailCard{
		thumbnail:           canvas.NewImageFromResource(nil),
		heading:             widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		meta:                widget.NewLabel(""),
		highlight:           canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
		playInBrowserButton: widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		addToPlaylistButton: widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
	}
	c.thumbnail.SetMinSize(fyne.NewSize(112, 63))
	c.highlight.Hide()
	c.heading.Truncation = fyne.TextTruncateEllipsis
	c.heading.Wrapping = fyne.TextWrapBreak
	c.meta.Truncation = fyne.TextTruncateEllipsis
	c.meta.Wrapping = fyne.TextWrapBreak
	c.playInBrowserButton.Importance = widget.LowImportance
	c.playInBrowserButton.OnTapped = func() { go playInBrowserCallback(c.result) }
	c.addToPlaylistButton.Importance = widget.LowImportance
	c.addToPlaylistButton.OnTapped = func() { go addToPlaylistCallback(c.result) }

	c.ExtendBaseWidget(&c)
	return &c
}

func (c *ThumbnailCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, nil, c.thumbnail, container.NewVBox(c.playInBrowserButton, c.addToPlaylistButton), container.NewVBox(c.heading, c.meta)),
		c.highlight,
	))
}

func (c *ThumbnailCard) MouseIn(*desktop.MouseEvent) {
	c.highlight.Show()
	c.Refresh()
}

func (c *ThumbnailCard) MouseOut() {
	c.highlight.Hide()
	c.Refresh()
}

func (c *ThumbnailCard) MouseMoved(*desktop.MouseEvent) {
	// Satisfy MouseMovement interface.
}

func (c *ThumbnailCard) Set(result scrapers.Result) {
	// Update thumbnail.
	scaledThumbnail, err := ScaleImageFromBytes(result.Thumbnail, int(c.thumbnail.MinSize().Width), int(c.thumbnail.MinSize().Height))
	if err != nil {
		slog.Error("failed to scale image", "error", err)
		return
	}

	c.thumbnail.Image = scaledThumbnail
	c.heading.SetText(result.Title)
	totalSeconds := int(result.Length.Round(time.Second).Seconds())
	mins := totalSeconds / 60
	secs := totalSeconds % 60
	c.meta.SetText(fmt.Sprintf("[%02d:%02d] %s • %s", mins, secs, result.ChannelTitle, result.Stats))
	c.result = result
	c.Refresh()
}
