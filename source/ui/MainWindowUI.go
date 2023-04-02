package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

const (
	mainWindowName = "Meowy Player"
)

var mainWindowSize fyne.Size
var mainWindowIcon fyne.Resource

func init() {
	const (
		mainWindowIconName = "icon.png"
	)

	mainWindowSize = fyne.NewSize(500, 650)

	var err error
	if mainWindowIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(mainWindowIconName)); err != nil {
		log.Fatal(err)
	}
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
