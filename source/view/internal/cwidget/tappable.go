package cwidget

import (
	"fyne.io/fyne/v2"
)

type TappableBase struct {
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func (t *TappableBase) Tapped(e *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped(e)
	}
}

func (t *TappableBase) TappedSecondary(e *fyne.PointEvent) {
	if t.OnTappedSecondary != nil {
		t.OnTappedSecondary(e)
	}
}
