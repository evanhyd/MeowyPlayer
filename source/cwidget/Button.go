package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/pattern"
)

type Button struct {
	widget.Button
	pattern.ZeroArgSubject
}

func NewButton(label string) *Button {
	button := &Button{}
	button.Text = label
	button.Importance = widget.LowImportance
	button.OnTapped = button.NotifyAll
	button.ExtendBaseWidget(button)
	return button
}

func NewButtonWithIcon(label string, icon fyne.Resource) *Button {
	button := NewButton(label)
	button.SetIcon(icon)
	return button
}

func (button *Button) SetOnTapped(onTapped func()) {
	button.OnTapped = func() {
		onTapped()
		button.NotifyAll()
	}
}
