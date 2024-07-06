package cwidget

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
	TappableComponent
	CursorableComponent
	text   *widget.RichText
	shadow *canvas.Rectangle
	music  model.Music
}

func NewMusicCardConstructor(onTapped func(model.Music), onTappedSecondary func(*fyne.PointEvent, model.Music)) func() *MusicCard {
	return func() *MusicCard {
		v := MusicCard{
			text:   widget.NewRichText(),
			shadow: canvas.NewRectangle(theme.HoverColor()),
		}
		v.OnTapped = func(*fyne.PointEvent) { onTapped(v.music) }
		v.OnTappedSecondary = func(e *fyne.PointEvent) { onTappedSecondary(e, v.music) }
		v.ExtendBaseWidget(&v)
		return &v
	}
}

func (v *MusicCard) CreateRenderer() fyne.WidgetRenderer {
	v.shadow.Hide()
	v.text.Wrapping = fyne.TextWrapWord
	v.text.Truncation = fyne.TextTruncateEllipsis
	return widget.NewSimpleRenderer(container.NewStack(v.text, v.shadow))
}

func (v *MusicCard) MouseIn(*desktop.MouseEvent) {
	v.shadow.Show()
	v.Refresh()
}

func (v *MusicCard) MouseOut() {
	v.shadow.Hide()
	v.Refresh()
}

func (v *MusicCard) MouseMoved(*desktop.MouseEvent) {
	//satisfy Hoverable interface
}

func (v *MusicCard) Notify(music model.Music) {
	length := music.Length().Round(time.Second)
	mins := length / time.Minute
	secs := (length - mins*time.Minute) / time.Second

	v.music = music
	v.text.Segments = v.text.Segments[:0]
	v.text.Segments = append(v.text.Segments, &widget.TextSegment{
		Text:  fmt.Sprintf("%02d:%02d | %v", mins, secs, music.Title()),
		Style: widget.RichTextStyleParagraph},
	)
	v.text.Refresh()
}
