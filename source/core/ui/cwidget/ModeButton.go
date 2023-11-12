package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ModeButton struct {
	widget.Button
	labels   []string
	icons    []fyne.Resource
	mode     int
	OnTapped func(int)
}

func newModeButton(labels []string, icons []fyne.Resource, onTapped func(int)) *ModeButton {
	button := &ModeButton{widget.Button{Importance: widget.LowImportance}, labels, icons, 0, onTapped}
	button.update()
	button.ExtendBaseWidget(button)
	return button
}

func (b *ModeButton) update() {
	if b.mode < len(b.labels) {
		b.SetText(b.labels[b.mode])
	}
	if b.mode < len(b.icons) {
		b.SetIcon(b.icons[b.mode])
	}
}

func (b *ModeButton) Tapped(*fyne.PointEvent) {
	maxLen := max(len(b.labels), len(b.icons))
	b.mode = (b.mode + 1) % maxLen
	b.update()
	b.OnTapped(b.mode)
}
