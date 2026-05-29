package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProfilePage struct {
	widget.BaseWidget
}

func newProfilePage() *ProfilePage {
	p := ProfilePage{}
	p.ExtendBaseWidget(&p)
	return &p
}

func (p *ProfilePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack())
}
