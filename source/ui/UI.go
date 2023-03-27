package ui

import (
	"log"
	"strings"

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
	menu.SetTabLocation(container.TabLocationLeading)
	menu.OnSelected = func(tab *container.TabItem) {
		tab.Content.Refresh()
	}

	meowyPlayerWindow.SetContent(menu)
	return meowyPlayerWindow
}

func createAblumTab() *container.TabItem {
	searchBar := cwidget.NewSearchBar()
	sortByNameButton := cwidget.NewButton("Name")
	sortByModifiedTimeButton := cwidget.NewButton("Date")

	scroll := cwidget.NewAlbumItemList(
		func() fyne.CanvasObject {
			card := cwidget.NewCard("", "", albumCoverIcon)
			title := widget.NewLabel("")
			return container.NewBorder(nil, nil, card, nil, title)
		},
		func(album player.Album, canvas fyne.CanvasObject) {
			//weak design, if the inner border style change, then this code would break easily
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
	sortByModifiedTimeButton.AddObserver(scroll.DateFilter())
	player.GetPlayerState().OnUpdateAllAlbumsAddObserver(scroll.ItemUpdater())
	sortByModifiedTimeButton.OnTapped()

	canvas := container.NewBorder(
		container.NewBorder(
			searchBar,
			nil,
			nil,
			nil,
			container.NewGridWithRows(1, sortByNameButton, sortByModifiedTimeButton),
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
	sortByModifiedTimeButton := cwidget.NewButton("Date")

	scroll := cwidget.NewItemList(
		func() fyne.CanvasObject {
			card := cwidget.NewCard("", "", albumCoverIcon)
			title := widget.NewLabel("")
			return container.NewBorder(nil, nil, card, nil, title)
		},
		func(album player.Album, canvas fyne.CanvasObject) {
			//weak design, if the inner border style change, then this code would break easily
			label := canvas.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(album.Description())

			card := canvas.(*fyne.Container).Objects[1].(*cwidget.Card)
			card.SetImage(album.CoverIcon())
			card.Image.SetMinSize(albumCoverIconSize)
		},
	)

	searchBar.SetOnChanged(func(text string) {
		lowerCaseText := strings.ToLower(text)
		scroll.SetFilter(func(album player.Album) bool {
			return strings.Contains(strings.ToLower(album.Title()), lowerCaseText)
		})
		scroll.ScrollToTop()
	})
	sortByNameButton.SetOnTapped(func() {
		scroll.SetSorter(func(album0, album1 player.Album) bool {
			return strings.Compare(strings.ToLower(album0.Title()), strings.ToLower(album1.Title())) < 0
		})
	})
	sortByModifiedTimeButton.SetOnTapped(func() {
		scroll.SetSorter(func(album0, album1 player.Album) bool {
			return album0.ModifiedTime().Compare(album1.ModifiedTime()) > 0
		})
	})

	canvas := container.NewBorder(
		container.NewBorder(
			searchBar,
			nil,
			nil,
			nil,
			container.NewGridWithRows(1, sortByNameButton, sortByModifiedTimeButton),
		),
		nil,
		nil,
		nil,
		scroll,
	)

	defer sortByModifiedTimeButton.OnTapped()
	return container.NewTabItemWithIcon(musicTabName, musicTabIcon, canvas)
}

func createAlbumAdderTab() *container.TabItem {
	return container.NewTabItemWithIcon(albumAdderTabName, albumAdderTabIcon, container.NewVBox(cwidget.NewButton("album adder")))
}

func createMusicAdderTab() *container.TabItem {
	return container.NewTabItemWithIcon(musicAdderTabName, musicAdderTabIcon, container.NewVBox(cwidget.NewButton("music adder")))
}
