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
	heading   *widget.RichText
	meta      *widget.RichText
	highlight *canvas.Rectangle
}

func NewThumbnailCard() *ThumbnailCard {
	c := ThumbnailCard{
		thumbnail: canvas.NewImageFromResource(nil),
		heading:   widget.NewRichTextWithText(""),
		meta:      widget.NewRichTextWithText(""),
		highlight: canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
	}
	c.ExtendBaseWidget(&c)
	return &c
}

func (c *ThumbnailCard) CreateRenderer() fyne.WidgetRenderer {
	c.thumbnail.SetMinSize(fyne.NewSize(112, 63))
	c.highlight.Hide()
	c.heading.Truncation = fyne.TextTruncateEllipsis
	c.heading.Wrapping = fyne.TextWrapBreak
	c.meta.Truncation = fyne.TextTruncateEllipsis
	c.meta.Wrapping = fyne.TextWrapBreak

	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, nil, c.thumbnail, nil, container.NewVBox(c.heading, c.meta)),
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
	scaledThumbnail, err := scaleImage(result.Thumbnail, int(c.thumbnail.MinSize().Width), int(c.thumbnail.MinSize().Height))
	if err != nil {
		fyne.LogError("Failed to decode or scale thumbnail", err)
		return
	}
	c.thumbnail.Image = scaledThumbnail

	// Update Heading.
	c.heading.Segments[0] = &widget.TextSegment{
		Style: widget.RichTextStyle{TextStyle: fyne.TextStyle{Bold: true}},
		Text:  result.Title,
	}

	// Update summary.
	totalSeconds := int(result.Length.Round(time.Second).Seconds())
	mins := totalSeconds / 60
	secs := totalSeconds % 60
	c.meta.Segments[0] = &widget.TextSegment{
		Text: fmt.Sprintf("[%02d:%02d] %s • %s", mins, secs, result.ChannelTitle, result.Stats),
	}

	c.Refresh()
}
