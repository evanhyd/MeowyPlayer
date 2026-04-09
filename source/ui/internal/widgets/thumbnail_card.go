package widgets

import (
	"fmt"
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
	thumbnail *canvas.Image
	summary   *widget.RichText
	highlight *canvas.Rectangle
}

func NewThumbnailCard() *ThumbnailCard {
	c := ThumbnailCard{
		thumbnail: canvas.NewImageFromResource(nil),
		summary:   widget.NewRichTextWithText(""),
		highlight: canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
	}
	c.ExtendBaseWidget(&c)
	return &c
}

func (c *ThumbnailCard) CreateRenderer() fyne.WidgetRenderer {
	c.thumbnail.SetMinSize(fyne.NewSize(64, 64))
	c.summary.Wrapping = fyne.TextWrapWord
	c.highlight.Hide()

	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, nil, c.thumbnail, nil, c.summary),
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
	scaledThumbnail, err := scaleImage(result.Thumbnail, 64, 64)
	if err != nil {
		fyne.LogError("Failed to decode or scale thumbnail", err)
		return
	}
	c.thumbnail.Image = scaledThumbnail

	// Update summary.
	totalSeconds := int(result.Length.Round(time.Second).Seconds())
	mins := totalSeconds / 60
	secs := totalSeconds % 60

	heading := widget.TextSegment{
		Style: widget.RichTextStyle{TextStyle: fyne.TextStyle{Bold: true}},
		Text:  fmt.Sprintf("[%02d:%02d] %s", mins, secs, result.Title),
	}
	meta := widget.TextSegment{
		Text: fmt.Sprintf("%s • %s", result.ChannelTitle, result.Stats),
	}
	c.summary.Segments = c.summary.Segments[:0]
	c.summary.Segments = append(c.summary.Segments, &heading, &meta)
	c.Refresh()
}
