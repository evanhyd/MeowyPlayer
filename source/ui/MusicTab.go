package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource/album"
	"meowyplayer.com/source/resource/texture"
)

type O struct {
}

func (o *O) Notify(album *player.Album) {
	fmt.Println("playing: " + album.Title)
}

func newMusicTab() *container.TabItem {
	const (
		musicTabName     = "Music"
		musicTabIconName = "music_tab.png"

		musicAdderLocalIconName        = "music_adder_local.png"
		musicAdderOnlineIconName       = "music_adder_online.png"
		musicAdderOnlineSearchIconName = "music_adder_online_search.png" //move to other place
	)

	o := O{}
	album.Get().Attach(&o)

	//search bar
	searchBar := widget.NewEntry()
	searchBar.OnChanged = func(title string) {
		//to do: set text filter
		//to do: scroll to top
	}

	//add local music
	musicAdderLocalButton := widget.NewButtonWithIcon("", texture.Get(musicAdderLocalIconName), func() {
		log.Println("add music from local")
		//to do
	})
	musicAdderLocalButton.Importance = widget.LowImportance

	//add online music
	musicAdderOnlineButton := widget.NewButtonWithIcon("", texture.Get(musicAdderOnlineIconName), func() {
		log.Println("add music from online")
		//to do
	})
	musicAdderOnlineButton.Importance = widget.LowImportance

	//sort by title button
	reverseTitle := false
	sortByTitleButton := widget.NewButton("Title", func() {
		reverseTitle = !reverseTitle
		//set sorter
	})
	sortByTitleButton.Importance = widget.LowImportance

	//sort by date button
	reverseDate := true
	sortByDateButton := widget.NewButton("Date", func() {
		reverseDate = !reverseDate
		//set sorter
	})
	sortByDateButton.Importance = widget.LowImportance
	sortByDateButton.OnTapped()

	border := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, sortByTitleButton, sortByDateButton),
			nil,
			container.NewGridWithRows(1, musicAdderLocalButton, musicAdderOnlineButton),
			searchBar,
		),
		nil,
		nil,
		nil,
		widget.NewLabel("haha"),
	)

	return container.NewTabItemWithIcon(musicTabName, texture.Get(musicTabIconName), border)
}
