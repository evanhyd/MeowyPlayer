package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Card struct {
	widget.Card
	OnTapped func()
}

// The last parameter is an image, not a content!!!
func NewCardWithImage(title string, subtitle string, content fyne.CanvasObject, image *canvas.Image) *Card {
	card := &Card{Card: widget.Card{Title: title, Subtitle: subtitle, Content: content, Image: image}}
	card.OnTapped = func() {}
	card.ExtendBaseWidget(card)
	return card
}

func (card *Card) Tapped(*fyne.PointEvent) {
	card.OnTapped()
}

func (card *Card) MouseMoved(*desktop.MouseEvent) {
	//MouseIn() MouseOut() interface
}

func (card *Card) MouseIn(*desktop.MouseEvent) {
	card.Image.Translucency = 0.2
	card.Refresh()
}

func (card *Card) MouseOut() {
	card.Image.Translucency = 0.0
	card.Refresh()
}
