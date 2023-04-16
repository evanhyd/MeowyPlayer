package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

const (
	musicTabName = "Music"
)

var musicTabIcon fyne.Resource

func init() {
	const (
		musicTabIconName = "music_tab.png"
	)

	var err error
	if musicTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicTabIconName)); err != nil {
		log.Fatal(err)
	}
}

func createMusicTab() *container.TabItem {
	searchBar := cwidget.NewSearchBar()
	sortByTitleButton := cwidget.NewButton("Title")
	sortByDateButton := cwidget.NewButton("Date")

	scroll := cwidget.NewMusicList(
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			button := cwidget.NewButton("<")
			return container.NewBorder(nil, nil, label, button)
		},
		func(music player.Music, canvas fyne.CanvasObject) {
			items := canvas.(*fyne.Container).Objects
			label := items[0].(*widget.Label)
			if label.Text != music.Description() {
				label.Text = music.Description()

				//update setting menu
				button := items[1].(*cwidget.Button)
				button.OnTapped = func() { DisplayErrorIfNotNil(player.RemoveMusicFromAlbum(player.GetState().Album(), music)) }

				canvas.Refresh()
			}
		},
	)

	searchBar.OnChanged = scroll.SetTitleFilter
	sortByTitleButton.OnTapped = scroll.SetTitleSorter
	sortByDateButton.OnTapped = scroll.SetDateSorter
	player.GetState().OnReadMusicFromDiskSubject().AddObserver(scroll)
	scroll.SetOnSelected(player.GetState().SetSelectedMusic)

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
		container.NewBorder(
			searchBar,
			nil,
			nil,
			nil,
			container.NewGridWithRows(1, sortByTitleButton, sortByDateButton),
		),
		nil,
		nil,
		nil,
		scroll,
	)
	return container.NewTabItemWithIcon(musicTabName, musicTabIcon, canvas)
}
