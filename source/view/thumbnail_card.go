package view

import (
	"fmt"
	"playground/browser"
	"playground/resource"
	"playground/view/internal/cwidget"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ThumbnailCard struct {
	widget.BaseWidget
	thumbnail *canvas.Image
	summary   *widget.RichText
	download  *widget.Button
	highlight *canvas.Rectangle
}

func newThumbnailCard() *ThumbnailCard {
	c := &ThumbnailCard{
		thumbnail: canvas.NewImageFromResource(theme.BrokenImageIcon()),
		summary:   widget.NewRichTextWithText(""),
		download:  cwidget.NewButtonWithIcon("", theme.DownloadIcon(), nil),
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *ThumbnailCard) CreateRenderer() fyne.WidgetRenderer {
	c.thumbnail.FillMode = canvas.ImageFillOriginal
	c.thumbnail.ScaleMode = canvas.ImageScaleFastest
	c.thumbnail.SetMinSize(resource.KThumbnailSize)
	c.summary.Wrapping = fyne.TextWrapWord
	return widget.NewSimpleRenderer(container.NewStack(container.NewBorder(nil, nil, c.thumbnail, c.download, c.summary), c.highlight))
}

func (v *ThumbnailCard) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.Refresh()
}

func (v *ThumbnailCard) MouseOut() {
	v.highlight.Hide()
	v.Refresh()
}

func (v *ThumbnailCard) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}

func (c *ThumbnailCard) Notify(result browser.Result) {
	length := result.Length.Round(time.Second)
	mins := length / time.Minute
	secs := (length - mins*time.Minute) / time.Second
	heading := &widget.TextSegment{
		Style: widget.RichTextStyle{TextStyle: fyne.TextStyle{Bold: true}},
		Text:  fmt.Sprintf("[%02d:%02d] %s", mins, secs, result.Title),
	}
	channel := &widget.TextSegment{Text: result.ChannelTitle}
	stats := &widget.TextSegment{Text: result.Stats}

	c.thumbnail.Resource = result.Thumbnail
	c.summary.Segments = c.summary.Segments[:0]
	c.summary.Segments = append(c.summary.Segments, heading, channel, stats)
	c.Refresh()
}
