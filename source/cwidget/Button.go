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
	button := &Button{Button: widget.Button{Text: label, Importance: widget.LowImportance}}
	button.ExtendBaseWidget(button)
	return button
}

func NewButtonWithIcon(label string, icon fyne.Resource) *Button {
	button := &Button{Button: widget.Button{Text: label, Importance: widget.LowImportance, Icon: icon}}
	button.ExtendBaseWidget(button)
	return button
}

func (button *Button) Tapped(*fyne.PointEvent) {
	button.Button.Tapped(nil)
	button.NotifyAll()
}

func (button *Button) OnTappedSubject() *pattern.ZeroArgSubject {
	return &button.ZeroArgSubject
}
