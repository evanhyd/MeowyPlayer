package mwidget

import (
	"fmt"
	"log/slog"
	"meowyplayer/scrapers"
	"meowyplayer/ui/internal/mutil"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = (*ThumbnailCard)(nil)

type ThumbnailCard struct {
	widget.BaseWidget
	highlight           *canvas.Rectangle
	thumbnail           *canvas.Image
	title               *widget.Label
	description         *widget.Label
	playInBrowserButton *widget.Button
	addToPlaylistButton *widget.Button
	result              scrapers.Result
}

func NewThumbnailCard(playInBrowserCallback func(scrapers.Result), addToPlaylistCallback func(scrapers.Result)) *ThumbnailCard {
	c := ThumbnailCard{
		highlight:           canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
		thumbnail:           canvas.NewImageFromResource(nil),
		title:               widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		description:         widget.NewLabel(""),
		playInBrowserButton: widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		addToPlaylistButton: widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
	}
	c.highlight.CornerRadius = 8.0
	c.highlight.Hide()
	c.thumbnail.SetMinSize(ThumbnailSize)
	c.thumbnail.CornerRadius = 8.0
	c.title.Truncation = fyne.TextTruncateEllipsis
	c.title.Wrapping = fyne.TextWrapBreak
	c.description.Truncation = fyne.TextTruncateEllipsis
	c.description.Wrapping = fyne.TextWrapBreak
	c.playInBrowserButton.Importance = widget.LowImportance
	c.playInBrowserButton.OnTapped = func() { go playInBrowserCallback(c.result) }
	c.addToPlaylistButton.Importance = widget.LowImportance
	c.addToPlaylistButton.OnTapped = func() { go addToPlaylistCallback(c.result) }

	c.ExtendBaseWidget(&c)
	return &c
}

func (c *ThumbnailCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, nil, c.thumbnail, container.NewVBox(c.playInBrowserButton, c.addToPlaylistButton), container.NewVBox(c.title, c.description)),
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

func (b *ThumbnailCard) Tapped(*fyne.PointEvent) {
	// Disable the yellow highlight from widget.List.
}

func (c *ThumbnailCard) Set(result scrapers.Result) {
	// Update thumbnail.
	scaledThumbnail, err := mutil.ScaleImageFromBytes(result.Thumbnail, c.thumbnail.MinSize())
	if err != nil {
		slog.Error("failed to scale image", "error", err)
		return
	}

	c.thumbnail.Image = scaledThumbnail
	c.title.SetText(result.Title)
	seconds := int64(result.Length.Round(time.Second).Seconds())
	c.description.SetText(fmt.Sprintf("[%s] %s • %s", mutil.SecondsToTime(seconds), result.ChannelTitle, result.Stats))
	c.result = result
	c.Refresh()
}
