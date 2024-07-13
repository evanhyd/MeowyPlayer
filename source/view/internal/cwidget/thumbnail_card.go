package cwidget

import (
	"fmt"
	"playground/browser"
	"playground/view/internal/resource"
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
	thumbnail   *canvas.Image
	summary     *widget.RichText
	instantPlay *widget.Button
	download    *widget.Button
	highlight   *canvas.Rectangle
	result      browser.Result
}

func NewThumbnailCardConstructor(
	onInstantPlay func(browser.Result),
	onDownload func(browser.Result)) func() *ThumbnailCard {
	return func() *ThumbnailCard {
		var c ThumbnailCard
		c = ThumbnailCard{
			thumbnail:   canvas.NewImageFromResource(theme.BrokenImageIcon()),
			summary:     widget.NewRichTextWithText(""),
			instantPlay: NewButtonIcon(theme.MediaPlayIcon(), func() { onInstantPlay(c.result) }),
			download:    NewButtonIcon(theme.DownloadIcon(), func() { onDownload(c.result) }),
			highlight:   canvas.NewRectangle(theme.HoverColor()),
		}
		c.ExtendBaseWidget(&c)
		return &c
	}
}

func (c *ThumbnailCard) CreateRenderer() fyne.WidgetRenderer {
	c.highlight.Hide()
	c.thumbnail.FillMode = canvas.ImageFillOriginal
	c.thumbnail.ScaleMode = canvas.ImageScaleFastest
	c.thumbnail.SetMinSize(resource.KThumbnailSize)
	c.summary.Wrapping = fyne.TextWrapWord
	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(
			nil,
			nil,
			c.thumbnail,
			nil,
			container.NewBorder(nil, container.NewHBox(c.instantPlay, c.download), nil, nil, c.summary),
		),
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

	c.result = result
	c.thumbnail.Resource = result.Thumbnail
	c.summary.Segments = c.summary.Segments[:0]
	c.summary.Segments = append(c.summary.Segments, heading, channel, stats)
	c.Refresh()
}
