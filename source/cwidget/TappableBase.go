package cwidget

import "fyne.io/fyne/v2"

type TappableBase struct {
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func (t *TappableBase) Tapped(event *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped(event)
	}
}

func (t *TappableBase) TappedSecondary(event *fyne.PointEvent) {
	if t.OnTappedSecondary != nil {
		t.OnTappedSecondary(event)
	}
}
