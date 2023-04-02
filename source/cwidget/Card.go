package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/pattern"
)

type Card struct {
	widget.Card
	pattern.ZeroArgSubject
	onTapped func()
}

// The last parameter is an image, not a content!!!
func NewCardWithImage(title string, subtitle string, image *canvas.Image) *Card {
	card := &Card{Card: widget.Card{Title: title, Subtitle: subtitle, Image: image}}
	card.onTapped = func() { card.NotifyAll() }
	card.ExtendBaseWidget(card)
	return card
}

func (card *Card) Tapped(*fyne.PointEvent) {
	card.onTapped()
}

func (card *Card) OnTapped() *pattern.ZeroArgSubject {
	return &card.ZeroArgSubject
}

func (card *Card) SetOnTapped(onTapped func()) {
	card.onTapped = func() {
		onTapped()
		card.NotifyAll()
	}
}
