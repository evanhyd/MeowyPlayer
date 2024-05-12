package view

import (
	"fmt"
	"playground/cwidget"
	"playground/model"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicCardProp struct {
	Music             model.Music
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

type MusicCard struct {
	widget.BaseWidget
	cwidget.TappableBase
	info      widget.TextSegment
	highlight *canvas.Rectangle
	isHovered bool
}

func newMusicCard() *MusicCard {
	v := &MusicCard{
		info:      widget.TextSegment{Style: widget.RichTextStyleParagraph},
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	v.highlight.Hide()
	v.ExtendBaseWidget(v)
	return v
}

func (v *MusicCard) CreateRenderer() fyne.WidgetRenderer {
	text := widget.NewRichText(&v.info)
	text.Wrapping = fyne.TextWrapWord
	text.Truncation = fyne.TextTruncateEllipsis
	return widget.NewSimpleRenderer(container.NewStack(text, v.highlight))
}

func (v *MusicCard) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.isHovered = true
	v.Refresh()
}

func (v *MusicCard) MouseOut() {
	v.highlight.Hide()
	v.isHovered = false
	v.Refresh()
}

func (v *MusicCard) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}

func (v *MusicCard) Cursor() desktop.Cursor {
	if v.isHovered {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (v *MusicCard) Notify(prop MusicCardProp) {
	length := prop.Music.Length().Round(time.Second)
	mins := length / time.Minute
	secs := (length - mins*time.Minute) / time.Second

	v.info.Text = fmt.Sprintf("%02d:%02d | %v", mins, secs, prop.Music.Title())
	v.Refresh()

	v.OnTapped = prop.OnTapped
	v.OnTappedSecondary = prop.OnTappedSecondary
}
