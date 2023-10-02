package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewButton(label string, tapped func()) *widget.Button {
	button := widget.NewButton(label, tapped)
	button.Importance = widget.LowImportance
	return button
}

func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(label, icon, tapped)
	button.Importance = widget.LowImportance
	return button
}
