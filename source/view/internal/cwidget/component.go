package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type TappableComponent struct {
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func (c *TappableComponent) Tapped(e *fyne.PointEvent) {
	if c.OnTapped != nil {
		c.OnTapped(e)
	}
}

func (c *TappableComponent) TappedSecondary(e *fyne.PointEvent) {
	if c.OnTappedSecondary != nil {
		c.OnTappedSecondary(e)
	}
}

type CursorableComponent struct {
	cursor desktop.Cursor
}

func (c *CursorableComponent) Cursor() desktop.Cursor {
	if c.cursor != nil {
		return c.cursor
	}
	return desktop.PointerCursor
}
