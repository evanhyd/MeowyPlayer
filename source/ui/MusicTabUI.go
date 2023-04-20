package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

const (
	musicTabName = "Music"
)

var musicTabIcon fyne.Resource
var musicAdderLocalIcon fyne.Resource
var musicAdderOnlineIcon fyne.Resource

func init() {
	const (
		musicTabIconName         = "music_tab.png"
		musicAdderLocalIconName  = "music_adder_local.png"
		musicAdderOnlineIconName = "music_adder_online.png"
	)

	var err error
	if musicTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicTabIconName)); err != nil {
		log.Fatal(err)
	}
	if musicAdderLocalIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicAdderLocalIconName)); err != nil {
		log.Fatal(err)
	}
	if musicAdderOnlineIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicAdderOnlineIconName)); err != nil {
		log.Fatal(err)
	}
}

func createMusicTab() *container.TabItem {
	musicAdderButton := cwidget.NewButtonWithIcon("", musicAdderLocalIcon)
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

	musicAdderButton.OnTapped = func() { createMusicAdderLocal().Show() }
	searchBar.OnChanged = scroll.SetTitleFilter
	sortByTitleButton.OnTapped = scroll.SetTitleSorter
	sortByDateButton.OnTapped = scroll.SetDateSorter
	player.GetState().OnUpdateMusics().AddCallback(scroll.Notify)
	scroll.SetOnSelected(func(music *player.Music) { player.UserSelectMusic(*music) })

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, sortByTitleButton, sortByDateButton),
			musicAdderButton,
			nil,
			searchBar,
		),
		nil,
		nil,
		nil,
		scroll,
	)
	return container.NewTabItemWithIcon(musicTabName, musicTabIcon, canvas)
}

func createMusicAdderLocal() *dialog.FileDialog {
	fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			DisplayErrorIfNotNil(err)
			return
		}
		if result != nil {
			DisplayErrorIfNotNil(player.AddMusicToAlbum(player.GetState().Album(), result.URI()))
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
	fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	fileOpenDialog.SetConfirmText("Add")
	return fileOpenDialog
}
