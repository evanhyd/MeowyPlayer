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
	info      *widget.Label
	highlight *canvas.Rectangle
	isHovered bool
}

func newMusicCard() *MusicCard {
	v := &MusicCard{
		info:      widget.NewLabel("?"),
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	v.highlight.Hide()
	v.ExtendBaseWidget(v)
	return v
}

func (v *MusicCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(v.info, v.highlight))
}

func (v *MusicCard) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.highlight.Refresh()
	v.isHovered = true
}

func (v *MusicCard) MouseOut() {
	v.highlight.Hide()
	v.highlight.Refresh()
	v.isHovered = false
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

	v.info.SetText(fmt.Sprintf("%02d:%02d | %v", mins, secs, prop.Music.Title()))
	v.OnTapped = prop.OnTapped
	v.OnTappedSecondary = prop.OnTappedSecondary
}
