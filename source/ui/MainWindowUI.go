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

var mainWindowIcon fyne.Resource

func init() {
	const (
		mainWindowIconName = "icon.png"
	)

	var err error
	if mainWindowIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(mainWindowIconName)); err != nil {
		log.Fatal(err)
	}
}

func NewMeowyPlayerWindow() fyne.Window {
	fyne.SetCurrentApp(app.New())
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())

	meowyPlayerWindow := fyne.CurrentApp().NewWindow(mainWindowName)
	meowyPlayerWindow.SetMaster()
	meowyPlayerWindow.Resize(resource.GetMainWindowSize())
	meowyPlayerWindow.SetIcon(mainWindowIcon)
	meowyPlayerWindow.CenterOnScreen()

	albumTab := createAblumTab()
	musicTab := createMusicTab()
	tabs := container.NewAppTabs(albumTab, musicTab)
	tabs.SetTabLocation(container.TabLocationLeading)

	//switch to the music tab after loaded music list
	player.GetState().OnFocusAlbumTab().AddCallback(func() { tabs.Select(albumTab) })
	player.GetState().OnFocusMusicTab().AddCallback(func() { tabs.Select(musicTab) })

	meowyPlayerWindow.SetContent(container.NewBorder(nil, createSeeker(), nil, nil, tabs))
	return meowyPlayerWindow
}
