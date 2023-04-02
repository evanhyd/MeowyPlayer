package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/resource"
)

const (
	musicAdderTabName = "Music Adder"
)

var musicAdderTabIcon fyne.Resource

func init() {
	const (
		musicAdderTabIconName = "music_adder_tab.png"
	)
	var err error
	if musicAdderTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicAdderTabIconName)); err != nil {
		log.Fatal(err)
	}
}

func createMusicAdderTab() *container.TabItem {
	return container.NewTabItemWithIcon(musicAdderTabName, musicAdderTabIcon, container.NewVBox(cwidget.NewButton("music adder")))
}
