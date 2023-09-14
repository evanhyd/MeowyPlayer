package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type CoverView struct {
	widget.Card
	OnTapped func()
}

// The last parameter is an image, not a content!!!
func NewCardWithImage(title string, subtitle string, content fyne.CanvasObject, image *canvas.Image) *CoverView {
	card := &CoverView{Card: widget.Card{Title: title, Subtitle: subtitle, Content: content, Image: image}}
	card.OnTapped = func() {}
	card.ExtendBaseWidget(card)
	return card
}

func (card *CoverView) Tapped(*fyne.PointEvent) {
	card.OnTapped()
}

func (card *CoverView) MouseIn(*desktop.MouseEvent) {
	card.Image.Translucency = 0.2
	card.Refresh()
}

func (card *CoverView) MouseOut() {
	card.Image.Translucency = 0.0
	card.Refresh()
}

func (card *CoverView) MouseMoved(*desktop.MouseEvent) {
	//MouseIn() MouseOut() interface
}
