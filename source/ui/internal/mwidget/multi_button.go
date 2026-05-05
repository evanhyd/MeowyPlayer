package mwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type MultiButton struct {
	widget.Button
	icons   []fyne.Resource
	actions []func()
	index   int
}

func NewMultiButton() *MultiButton {
	b := MultiButton{}
	b.Importance = widget.LowImportance
	b.Button.OnTapped = func() {
		b.Select((b.index + 1) % len(b.icons))
	}
	b.ExtendBaseWidget(&b)
	return &b
}

func (b *MultiButton) Select(index int) {
	b.index = index
	b.SetIcon(b.icons[index])
	b.actions[index]()
}

func (b *MultiButton) Add(icon fyne.Resource, action func()) {
	b.icons = append(b.icons, icon)
	b.actions = append(b.actions, action)
}
