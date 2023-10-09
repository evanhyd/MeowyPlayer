package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DropDown struct {
	widget.BaseWidget
	currentOption fyne.CanvasObject
	options       []*fyne.MenuItem
	menu          *widget.Menu
}

func NewDropDown() *DropDown {
	dropDown := &DropDown{}
	dropDown.ExtendBaseWidget(dropDown)
	return dropDown
}

func (d *DropDown) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.menu)
}

func (d *DropDown) Add(label string, object fyne.CanvasObject) {
	d.options = append(d.options, fyne.NewMenuItem(label, func() {
		d.currentOption = object
		d.Refresh()
	}))
}

func (d *DropDown) Tapped(*fyne.PointEvent) {
	// fyne.NewMenu()
	// widget.NewMenu()
	// widget.NewPopUpMenu()
	// rename := fyne.NewMenuItem("Rename", makeRenameDialog(album))
	// cover := fyne.NewMenuItem("Cover", makeCoverDialog(album))
	// delete := fyne.NewMenuItem("Delete", makeDeleteAlbumDialog(album))
	// return widget.NewPopUpMenu(fyne.NewMenu("", rename, cover, delete), canvas)
}
