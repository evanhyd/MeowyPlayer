package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewButtonWithIcon(label string, icon fyne.Resource, tapped func()) *widget.Button {
	button := widget.NewButtonWithIcon(label, icon, tapped)
	button.Importance = widget.LowImportance
	return button
}
