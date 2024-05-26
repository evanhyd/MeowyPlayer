package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DropDown struct {
	widget.BaseWidget
	selected *widget.Button
	menu     *fyne.Menu
}

func NewDropDown(icon fyne.Resource) *DropDown {
	d := &DropDown{
		selected: NewButtonWithIcon("", icon, nil),
		menu:     fyne.NewMenu(""),
	}

	d.selected.OnTapped = func() {
		canvas := fyne.CurrentApp().Driver().CanvasForObject(d)
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(d)
		widget.ShowPopUpMenuAtPosition(d.menu, canvas, position)
	}

	d.ExtendBaseWidget(d)
	return d
}

func (d *DropDown) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.selected)
}

func (d *DropDown) Add(title string, icon fyne.Resource, onSelected func()) {
	item := &fyne.MenuItem{
		Label: title,
		Icon:  icon,
		Action: func() {
			d.selected.SetIcon(icon)
			onSelected()
		},
	}
	d.menu.Items = append(d.menu.Items, item)
}

func (d *DropDown) Select(index int) {
	d.menu.Items[index].Action()
}
