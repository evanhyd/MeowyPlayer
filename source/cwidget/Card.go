package cwidget

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Card struct {
	widget.Card
}

func NewCard(title string, subtitle string, image *canvas.Image) *Card {
	card := &Card{}
	card.Title = title
	card.Subtitle = subtitle
	card.Image = image
	card.ExtendBaseWidget(card)
	return card
}
