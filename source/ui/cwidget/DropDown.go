package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DropDown struct {
	widget.BaseWidget
	sign      *Sign
	menu      *fyne.Menu
	highlight *canvas.Rectangle
}

func NewDropDown(title string, icon fyne.Resource) *DropDown {
	dropDown := &DropDown{
		sign:      NewSign(title, icon),
		menu:      fyne.NewMenu(""),
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	dropDown.highlight.Hide()
	dropDown.ExtendBaseWidget(dropDown)
	return dropDown
}

func (d *DropDown) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewMax(d.sign, d.highlight))
}

func (d *DropDown) Add(title string, icon fyne.Resource, onSelected func()) {
	item := fyne.NewMenuItem(title, func() {
		d.sign.title.SetText(title)
		d.sign.icon.SetResource(icon)
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

func (d *DropDown) MouseIn(*desktop.MouseEvent) {
	d.highlight.Show()
	d.Refresh()
}

func (d *DropDown) MouseOut() {
	d.highlight.Hide()
	d.Refresh()
}

func (d *DropDown) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
