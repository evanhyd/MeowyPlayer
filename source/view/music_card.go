package view

import (
	"fmt"
	"playground/model"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicCard struct {
	widget.BaseWidget
	info      widget.TextSegment
	highlight *canvas.Rectangle
}

func newMusicCard() *MusicCard {
	v := &MusicCard{
		info:      widget.TextSegment{Style: widget.RichTextStyleParagraph},
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	v.ExtendBaseWidget(v)
	return v
}

func (v *MusicCard) CreateRenderer() fyne.WidgetRenderer {
	v.highlight.Hide()
	text := widget.NewRichText(&v.info)
	text.Wrapping = fyne.TextWrapWord
	text.Truncation = fyne.TextTruncateEllipsis
	return widget.NewSimpleRenderer(container.NewStack(text, v.highlight))
}

func (v *MusicCard) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.Refresh()
}

func (v *MusicCard) MouseOut() {
	v.highlight.Hide()
	v.Refresh()
}

func (v *MusicCard) MouseMoved(*desktop.MouseEvent) {
	//satisfy Hoverable interface
}

func (v *MusicCard) Notify(music model.Music) {
	length := music.Length().Round(time.Second)
	mins := length / time.Minute
	secs := (length - mins*time.Minute) / time.Second

	v.info.Text = fmt.Sprintf("%02d:%02d | %v", mins, secs, music.Title())
	v.Refresh()
}
