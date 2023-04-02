package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

const (
	mainWindowName    = "Meowy Player"
	albumTabName      = "Album"
	musicTabName      = "Music"
	albumAdderTabName = "Album Adder"
	musicAdderTabName = "Music Adder"
)

var mainWindowSize fyne.Size
var albumCoverIconSize fyne.Size
var mainWindowIcon fyne.Resource
var albumTabIcon fyne.Resource
var musicTabIcon fyne.Resource
var albumAdderTabIcon fyne.Resource
var musicAdderTabIcon fyne.Resource
var albumCoverIcon *canvas.Image

func init() {
	const (
		mainWindowIconName    = "icon.png"
		albumTabIconName      = "album_tab.png"
		musicTabIconName      = "music_tab.png"
		albumAdderTabIconName = "album_adder_tab.png"
		musicAdderTabIconName = "music_adder_tab.png"
		albumCoverIconName    = "album_cover.png"
	)

	mainWindowSize = fyne.NewSize(500, 650)
	albumCoverIconSize = fyne.NewSize(128.0, 128.0)

	var err error
	if mainWindowIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(mainWindowIconName)); err != nil {
		log.Fatal(err)
	}
	if albumTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumTabIconName)); err != nil {
		log.Fatal(err)
	}
	if musicTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicTabIconName)); err != nil {
		log.Fatal(err)
	}
	if albumAdderTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumAdderTabIconName)); err != nil {
		log.Fatal(err)
	}
	if musicAdderTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(musicAdderTabIconName)); err != nil {
		log.Fatal(err)
	}
	albumCoverIcon = canvas.NewImageFromFile(resource.GetResourcePath(albumCoverIconName))
	albumCoverIcon.SetMinSize(albumCoverIconSize)
}

func NewMeowyPlayerWindow() fyne.Window {
	fyne.SetCurrentApp(app.New())
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())

	meowyPlayerWindow := fyne.CurrentApp().NewWindow(mainWindowName)
	meowyPlayerWindow.Resize(mainWindowSize)
	meowyPlayerWindow.SetIcon(mainWindowIcon)
	meowyPlayerWindow.CenterOnScreen()

	albumTab := createAblumTab()
	musicTab := createMusicTab()
	albumAdderTab := createAlbumAdderTab()
	musicAdderTab := createMusicAdderTab()
	menu := container.NewAppTabs(albumTab, musicTab, albumAdderTab, musicAdderTab)
	menu.SetTabLocation(container.TabLocationLeading)

	//switch to the music tab after loaded music list
	player.GetState().OnSelectAlbum().AddCallback(func(player.Album) { menu.SelectIndex(1) })

	meowyPlayerWindow.SetContent(container.NewBorder(nil, createSeeker(), nil, nil, menu))
	return meowyPlayerWindow
}

func createAblumTab() *container.TabItem {
	searchBar := cwidget.NewSearchBar()
	sortByTitleButton := cwidget.NewButton("Title")
	sortByDateButton := cwidget.NewButton("Date")

	scroll := cwidget.NewAlbumItemList(
		func() fyne.CanvasObject {
			card := widget.NewCard("", "", nil)
			card.SetImage(albumCoverIcon)
			title := widget.NewLabel("")
			return container.NewBorder(nil, nil, card, nil, title)
		},
		func(album player.Album, canvas fyne.CanvasObject) {
			//not a solid design. If the inner border style change, then this code would break
			label := canvas.(*fyne.Container).Objects[0].(*widget.Label)
			if label.Text != album.Description() {
				label.SetText(album.Description())
			}

			card := canvas.(*fyne.Container).Objects[1].(*widget.Card)
			if card.Image != album.CoverIcon() {
				card.SetImage(album.CoverIcon())
				card.Image.SetMinSize(albumCoverIconSize)
			}
		},
	)

	searchBar.SetOnChanged(scroll.SetTitleFilter)
	sortByTitleButton.SetOnTapped(scroll.SetTitleSorter)
	sortByDateButton.SetOnTapped(scroll.SetDateSorter)
	player.GetState().OnReadAlbumsFromDisk().AddObserver(scroll)
	scroll.SetOnSelected(player.GetState().SetSelectedAlbum)

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
	return container.NewTabItemWithIcon(albumTabName, albumTabIcon, canvas)
}

func createMusicTab() *container.TabItem {
	searchBar := cwidget.NewSearchBar()
	sortByTitleButton := cwidget.NewButton("Title")
	sortByDateButton := cwidget.NewButton("Date")

	scroll := cwidget.NewMusicItemList(
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(music player.Music, canvas fyne.CanvasObject) {
			label := canvas.(*widget.Label)
			if label.Text != music.Description() {
				label.SetText(music.Description())
			}
		},
	)

	searchBar.SetOnChanged(scroll.SetTitleFilter)
	sortByTitleButton.SetOnTapped(scroll.SetTitleSorter)
	sortByDateButton.SetOnTapped(scroll.SetDateSorter)
	player.GetState().OnReadMusicFromDisk().AddObserver(scroll)
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

func createAlbumAdderTab() *container.TabItem {
	return container.NewTabItemWithIcon(albumAdderTabName, albumAdderTabIcon, container.NewVBox(cwidget.NewButton("album adder")))
}

func createMusicAdderTab() *container.TabItem {
	return container.NewTabItemWithIcon(musicAdderTabName, musicAdderTabIcon, container.NewVBox(cwidget.NewButton("music adder")))
}

func createSeeker() *fyne.Container {
	albumView := cwidget.NewCardWithImage("", "", albumCoverIcon)
	player.GetState().OnSelectMusic().AddCallback(func(album player.Album, _ []player.Music, _ player.Music) {
		albumView.SetImage(album.CoverIcon())
		albumView.SetOnTapped(func() { player.GetState().SetSelectedAlbum(album) })
	})

	title := widget.NewLabel("label")
	player.GetState().OnSelectMusic().AddCallback(func(_ player.Album, _ []player.Music, music player.Music) {
		title.SetText(music.Title())
	})

	progressLabel := widget.NewLabel("00:00")
	progress := widget.NewSlider(0.0, 1.0)
	progress.Step = 1.0 / 1000.0

	previousButton := cwidget.NewButton(" << ")
	previousButton.SetOnTapped(player.GetPlayer().NextMusic)

	playButton := cwidget.NewButton(" O ")
	nextButton := cwidget.NewButton(" >> ")
	playModeButton := cwidget.NewButton("play mode")
	volume := widget.NewSlider(0.0, 1.0)
	volume.Step = 1.0 / 100.0

	return container.NewBorder(
		nil,
		nil,
		albumView,
		nil,
		container.NewBorder(
			title,
			container.NewHBox(layout.NewSpacer(), previousButton, playButton, nextButton, playModeButton, volume, layout.NewSpacer()),
			nil,
			nil,
			container.NewBorder(nil, nil, progressLabel, nil, progress),
		),
	)
}
