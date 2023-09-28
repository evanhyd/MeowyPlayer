package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/utility"
)

func NewMainWindow() fyne.Window {
	const windowTitle = "MeowyPlayer"
	windowSize := fyne.NewSize(770.0, 650.0)

	fyne.SetCurrentApp(app.NewWithID(windowTitle))
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())

	//set up windows orientation
	window := fyne.CurrentApp().NewWindow(windowTitle)
	window.SetMaster()
	window.SetIcon(resource.WindowIcon())
	window.SetCloseIntercept(window.Hide)
	window.Resize(windowSize)
	window.CenterOnScreen()

	//create system tray
	if desktop, ok := fyne.CurrentApp().(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
		desktop.SetSystemTrayIcon(resource.WindowIcon())
	}

	//set up item tabs
	albumTab := newAlbumTab()
	musicTab := newMusicTab()
	tabs := container.NewAppTabs(albumTab, musicTab)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.DisableItem(musicTab)
	client.GetCurrentAlbum().Attach(utility.MakeCallback(func(*player.Album) {
		tabs.EnableItem(musicTab)
		tabs.Select(musicTab)
	}))

	window.SetContent(container.NewBorder(nil, newController(), nil, nil, tabs))
	return window
}

func getMainWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func showErrorIfAny(err error) {
	if err != nil {
		log.Printf("gui error: %v\n", err)
		dialog.ShowError(err, getMainWindow())
	}
}
