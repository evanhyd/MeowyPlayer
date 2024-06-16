package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
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

type ToolbarButton struct {
	*widget.Button
}

func (t *ToolbarButton) ToolbarObject() fyne.CanvasObject {
	return t
}

func NewToolbarButton(label string, icon fyne.Resource, tapped func()) *ToolbarButton {
	return &ToolbarButton{NewButtonWithIcon(label, icon, tapped)}
}

func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(label, icon, tapped)
	button.Importance = widget.LowImportance
	return button
}

func NewMenuItemWithIcon(label string, icon fyne.Resource, action func()) *fyne.MenuItem {
	return &fyne.MenuItem{Label: label, Icon: icon, Action: action}
}
