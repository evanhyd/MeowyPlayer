package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	return container.NewTabItemWithIcon(albumAdderTabName, albumAdderTabIcon, container.NewVBox(cwidget.NewButton("album adder")))
}
