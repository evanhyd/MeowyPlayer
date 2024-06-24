package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type toolbarWidget struct {
	fyne.Widget
}

func (t *toolbarWidget) ToolbarObject() fyne.CanvasObject {
	return t.Widget
}

func NewButton(label string, icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(label, icon, tapped)
	button.Importance = widget.LowImportance
	return button
}

func NewTappableIcon(icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon("", icon, tapped)
	button.Importance = widget.LowImportance
	return button
}

func NewMenuItem(label string, icon fyne.Resource, action func()) *fyne.MenuItem {
	return &fyne.MenuItem{Label: label, Icon: icon, Action: action}
}
