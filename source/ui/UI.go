package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

var mainWindowSize fyne.Size
var albumCoverIconSize fyne.Size

var mainWindowIcon fyne.Resource
var albumTabIcon fyne.Resource
var musicTabIcon fyne.Resource
var albumAdderTabIcon fyne.Resource
var musicAdderTabIcon fyne.Resource

var albumCoverIcon *canvas.Image

const (
	mainWindowName    = "Meowy Player"
	albumTabName      = "Album"
	musicTabName      = "Music"
	albumAdderTabName = "Album Adder"
	musicAdderTabName = "Music Adder"

	mainWindowIconName    = "icon.png"
	albumTabIconName      = "album_tab.png"
	musicTabIconName      = "music_tab.png"
	albumAdderTabIconName = "album_adder_tab.png"
	musicAdderTabIconName = "music_adder_tab.png"
	albumCoverIconName    = "album_cover.png"
)

type MenuController struct {
	menu *container.AppTabs
}

func (menuController *MenuController) Notify(player.Album, []player.Music) {
	menuController.menu.SelectIndex(1)
}

func init() {
	mainWindowSize = fyne.NewSize(733.3747416, 733.3747416/1.618)
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
	menu.OnSelected = func(tab *container.TabItem) { tab.Content.Refresh() }
	menu.SetTabLocation(container.TabLocationLeading)

	//switch to the music tab after loaded music list
	player.GetState().OnSelectAlbum().AddObserver(&MenuController{menu})

	meowyPlayerWindow.SetContent(menu)
	return meowyPlayerWindow
}

func createAblumTab() *container.TabItem {
	searchBar := cwidget.NewSearchBar()
	sortByNameButton := cwidget.NewButton("Name")
	sortByDateButton := cwidget.NewButton("Date")

	scroll := cwidget.NewAlbumItemList(
		func() fyne.CanvasObject {
			card := cwidget.NewCard("", "", albumCoverIcon)
			title := widget.NewLabel("")
			return container.NewBorder(nil, nil, card, nil, title)
		},
		func(album player.Album, canvas fyne.CanvasObject) {
			//not a solid design. If the inner border style change, then this code would break
			label := canvas.(*fyne.Container).Objects[0].(*widget.Label)
			if label.Text != album.Description() {
				label.SetText(album.Description())
			}

			card := canvas.(*fyne.Container).Objects[1].(*cwidget.Card)
			if card.Image != album.CoverIcon() {
				card.SetImage(album.CoverIcon())
				card.Image.SetMinSize(albumCoverIconSize)
			}
		},
	)

	searchBar.AddObserver(scroll.NameFilter())
	sortByNameButton.AddObserver(scroll.NameSorter())
	sortByDateButton.AddObserver(scroll.DateFilter())
	player.GetState().OnReadAlbums().AddObserver(scroll.ItemUpdater())
	scroll.OnSelected().AddObserver(player.GetState().Info())

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
		container.NewBorder(
			searchBar,
			nil,
			nil,
			nil,
			container.NewGridWithRows(1, sortByNameButton, sortByDateButton),
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
	sortByNameButton := cwidget.NewButton("Name")
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

	searchBar.AddObserver(scroll.NameFilter())
	sortByNameButton.AddObserver(scroll.NameSorter())
	sortByDateButton.AddObserver(scroll.DateFilter())
	player.GetState().OnSelectAlbum().AddObserver(scroll.ItemUpdater())
	// scroll.OnSelected().AddObserver(player.GetState().SelectedAlbumUpdater())

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
		container.NewBorder(
			searchBar,
			nil,
			nil,
			nil,
			container.NewGridWithRows(1, sortByNameButton, sortByDateButton),
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
