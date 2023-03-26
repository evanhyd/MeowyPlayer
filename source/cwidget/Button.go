package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ButtonObserver interface {
	Notify()
}

type Button struct {
	widget.Button
	observers []ButtonObserver
}

func NewButton(label string) *Button {
	button := &Button{}
	button.Text = label
	button.Importance = widget.LowImportance
	button.OnTapped = button.NotifyObservers
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
		button.NotifyObservers()
	}
}

func (button *Button) AddObserver(observer ButtonObserver) {
	button.observers = append(button.observers, observer)
}

func (button *Button) NotifyObservers() {
	for _, observer := range button.observers {
		observer.Notify()
	}
}
