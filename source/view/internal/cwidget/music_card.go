package cwidget

import (
	"fmt"
	"playground/model"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type MusicCard struct {
	widget.BaseWidget
	TappableComponent
	CursorableComponent
	text  *widget.RichText
	music model.Music
}

func NewMusicCardConstructor(onTapped func(model.Music), onTappedSecondary func(*fyne.PointEvent, model.Music)) func() *MusicCard {
	return func() *MusicCard {
		v := MusicCard{text: widget.NewRichText()}
		v.OnTapped = func(*fyne.PointEvent) { onTapped(v.music) }
		v.OnTappedSecondary = func(e *fyne.PointEvent) { onTappedSecondary(e, v.music) }
		v.ExtendBaseWidget(&v)
		return &v
	}
}

func (v *MusicCard) CreateRenderer() fyne.WidgetRenderer {
	v.text.Wrapping = fyne.TextWrapWord
	v.text.Truncation = fyne.TextTruncateEllipsis
	return widget.NewSimpleRenderer(v.text)
}

func (v *MusicCard) Notify(music model.Music) {
	mins := music.Length() / time.Minute
	secs := (music.Length() - mins*time.Minute) / time.Second

	v.music = music
	v.text.Segments = v.text.Segments[:0]
	v.text.Segments = append(v.text.Segments, &widget.TextSegment{
		Text:  fmt.Sprintf("%02d:%02d | %v", mins, secs, music.Title()),
		Style: widget.RichTextStyleParagraph},
	)
	v.text.Refresh()
}
