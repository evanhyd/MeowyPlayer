package cwidget

import "fyne.io/fyne/v2"

type tappableBase struct {
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func (t *tappableBase) Tapped(event *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped(event)
	}
}

func (t *tappableBase) TappedSecondary(event *fyne.PointEvent) {
	if t.OnTappedSecondary != nil {
		t.OnTappedSecondary(event)
	}
}
