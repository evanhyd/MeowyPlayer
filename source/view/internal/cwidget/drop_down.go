package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DropDown struct {
	widget.BaseWidget
	selected *widget.Button
	menu     *fyne.Menu
}

func NewDropDown() *DropDown {
	var d DropDown
	d = DropDown{
		selected: NewTappableIcon(nil, d.showMenu),
		menu:     fyne.NewMenu(""),
	}
	d.ExtendBaseWidget(&d)
	return &d
}

func (d *DropDown) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.selected)
}

func (d *DropDown) showMenu() {
	canvas := fyne.CurrentApp().Driver().CanvasForObject(d)
	position := fyne.CurrentApp().Driver().AbsolutePositionForObject(d)
	position.Y += d.Size().Height - theme.InputBorderSize()
	widget.ShowPopUpMenuAtPosition(d.menu, canvas, position)
}

func (d *DropDown) Add(item *fyne.MenuItem) {
	action := item.Action
	item.Action = func() {
		d.selected.SetIcon(item.Icon)
		action()
	}
	d.menu.Items = append(d.menu.Items, item)

	if len(d.menu.Items) == 1 {
		item.Action()
	}
}
