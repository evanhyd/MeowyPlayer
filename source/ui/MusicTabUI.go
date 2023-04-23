package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/scraper"
)

const (
	musicTabName = "Music"
)

var musicTabIcon fyne.Resource
var musicAdderLocalIcon fyne.Resource
var musicAdderOnlineIcon fyne.Resource
var musicAdderOnlineSearchIcon fyne.Resource

func init() {
	const (
		musicTabIconName               = "music_tab.png"
		musicAdderLocalIconName        = "music_adder_local.png"
		musicAdderOnlineIconName       = "music_adder_online.png"
		musicAdderOnlineSearchIconName = "music_adder_online_search.png"
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
	if musicAdderOnlineSearchIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicAdderOnlineSearchIconName)); err != nil {
		log.Fatal(err)
	}
}

func createMusicTab() *container.TabItem {
	musicAdderLocalButton := cwidget.NewButtonWithIcon("", musicAdderLocalIcon)
	musicAdderOnlineButton := cwidget.NewButtonWithIcon("", musicAdderOnlineIcon)
	searchBar := widget.NewEntry()
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
				button := items[1].(*cwidget.Button)
				button.OnTapped = func() { DisplayErrorIfAny(player.RemoveMusicFromAlbum(player.GetState().Album(), music)) }
				canvas.Refresh()
			}
		},
	)

	musicAdderLocalButton.OnTapped = func() { createAddLocalDialog().Show() }
	musicAdderOnlineButton.OnTapped = func() { createAddOnlineDialog().Show() }
	searchBar.OnChanged = scroll.SetTitleFilter
	sortByTitleButton.OnTapped = scroll.SetTitleSorter
	sortByDateButton.OnTapped = scroll.SetDateSorter
	player.GetState().OnUpdateMusics().AddCallback(scroll.SetItems)
	scroll.SetOnSelected(func(music *player.Music) { player.UserSelectMusic(*music) })

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
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
		scroll,
	)
	return container.NewTabItemWithIcon(musicTabName, musicTabIcon, canvas)
}

func createAddLocalDialog() *dialog.FileDialog {
	fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			DisplayErrorIfAny(err)
			return
		}
		if result != nil {
			DisplayErrorIfAny(player.AddMusicToAlbum(player.GetState().Album(), result.URI().Path(), result.URI().Name()))
		}
	}, player.GetMainWindow())
	fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	fileOpenDialog.SetConfirmText("Add")
	return fileOpenDialog
}

func createAddOnlineDialog() dialog.Dialog {
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder(">>")
	searchButton := cwidget.NewButtonWithIcon("", musicAdderOnlineSearchIcon)

	scroll := cwidget.NewList(
		func() fyne.CanvasObject {
			card := widget.NewCard("", "", nil)
			card.Image = canvas.NewImageFromResource(defaultIcon)
			card.Image.SetMinSize(resource.GetThumbnailIconSize())

			videoTitle := widget.NewLabel("")
			videoTitle.TextStyle = fyne.TextStyle{Bold: true, Monospace: true, Symbol: true}

			videoInfo := widget.NewLabel("")
			description := widget.NewLabel("")

			return container.NewBorder(
				nil,
				nil,
				card,
				nil,
				container.NewGridWithRows(3, videoTitle, videoInfo, description),
			)
		},

		func(result scraper.ClipzagResult, canvas fyne.CanvasObject) {
			borderItems := canvas.(*fyne.Container).Objects
			gridItems := borderItems[0].(*fyne.Container).Objects

			videoTitle := gridItems[0].(*widget.Label)
			if videoTitle.Text != result.VideoTitle() {
				card := borderItems[1].(*widget.Card)
				card.Image = result.Thumbnail()

				videoTitle.Text = result.VideoTitle()

				videoInfo := gridItems[1].(*widget.Label)
				videoInfo.Text = result.ChannelTitle() + " | " + result.Stats()

				description := gridItems[2].(*widget.Label)
				description.Text = result.Description()

				canvas.Refresh()
			}
		},
	)

	searchButton.OnTapped = func() {
		result, err := scraper.GetSearchResult(searchBar.Text)
		if err != nil {
			DisplayErrorIfAny(err)
			return
		}
		scroll.SetItems(result)
	}

	searchBar.OnSubmitted = func(title string) {
		searchButton.OnTapped()
	}

	scroll.SetOnSelected(func(result *scraper.ClipzagResult) {
		DisplayErrorIfAny(scraper.AddMusicToRepository(result.VideoID(), player.GetState().Album(), result.VideoTitle()))
	})

	onlineBrowserDialog := dialog.NewCustom("", "( X )", container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			nil,
			searchButton,
			searchBar,
		),
		nil,
		nil,
		nil,
		scroll,
	), player.GetMainWindow())
	onlineBrowserDialog.Resize(resource.GetMusicAddOnlineDialogSize())
	return onlineBrowserDialog
}
