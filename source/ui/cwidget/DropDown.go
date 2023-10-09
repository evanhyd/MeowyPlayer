package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DropDown struct {
	Sign
	menu *fyne.Menu
}

func NewDropDown(title string, icon fyne.Resource) *DropDown {
	dropDown := &DropDown{Sign{title: widget.NewLabel(title), icon: widget.NewIcon(icon)}, fyne.NewMenu("")}
	dropDown.ExtendBaseWidget(dropDown)
	return dropDown
}

func (d *DropDown) Add(title string, icon fyne.Resource, onSelected func()) {
	item := fyne.NewMenuItem(title, func() {
		d.title.SetText(title)
		d.icon.SetResource(icon)
		onSelected()
	})
	item.Icon = icon
	d.menu.Items = append(d.menu.Items, item)
}

func (d *DropDown) Select(index int) {
	d.menu.Items[index].Action()
}

func (d *DropDown) Tapped(event *fyne.PointEvent) {
	canvas := fyne.CurrentApp().Driver().CanvasForObject(d)
	position := fyne.CurrentApp().Driver().AbsolutePositionForObject(d)
	widget.ShowPopUpMenuAtPosition(d.menu, canvas, position)
}
