package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/resource"
)

const (
	albumAdderTabName = "Album Adder"
)

var albumAdderTabIcon fyne.Resource

func init() {
	const (
		albumAdderTabIconName = "album_adder_tab.png"
	)
	var err error
	if albumAdderTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumAdderTabIconName)); err != nil {
		log.Fatal(err)
	}
}

func createAlbumAdderTab() *container.TabItem {
	but := cwidget.NewButton("button")
	box := container.NewVBox(but)

	but.SetOnTapped(func() {
		menuItem1 := fyne.NewMenuItem("A", func() {})
		menuItem2 := fyne.NewMenuItem("B", func() {})
		menuItem3 := fyne.NewMenuItem("C", func() {})
		menu := fyne.NewMenu("MENUUU", menuItem1, menuItem2, menuItem3)
		popUpMenu := widget.NewPopUpMenu(menu, fyne.CurrentApp().Driver().CanvasForObject(box))
		popUpMenu.ShowAtPosition(fyne.NewPos(200.0, 200.0))
	})

	return container.NewTabItemWithIcon(albumAdderTabName, albumAdderTabIcon, box)
}
