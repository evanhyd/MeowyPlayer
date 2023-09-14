package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/utility"
)

func NewMainWindow() fyne.Window {
	const (
		iconName    = "icon.ico"
		windowTitle = "MeowyPlayer"
	)
	windowSize := fyne.NewSize(490.0, 650.0)
	icon := resource.GetAsset(iconName)

	fyne.SetCurrentApp(app.NewWithID(windowTitle))
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())

	//set up windows orientation
	window := fyne.CurrentApp().NewWindow(windowTitle)
	window.SetMaster()
	window.SetIcon(icon)
	window.SetCloseIntercept(window.Hide)
	window.Resize(windowSize)
	window.CenterOnScreen()

	//create system tray
	if desktop, ok := fyne.CurrentApp().(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
		desktop.SetSystemTrayIcon(icon)
	}

	//set up item tabs
	albumTab := newAlbumTab()
	musicTab := newMusicTab()
	tabs := container.NewAppTabs(albumTab, musicTab)
	tabs.SetTabLocation(container.TabLocationLeading)
	manager.GetCurrentAlbum().Attach(utility.MakeCallback(func(_ *player.Album) { tabs.Select(musicTab) }))

	controller := newController()
	window.SetContent(container.NewBorder(nil, controller, nil, nil, tabs))
	window.Canvas().Scale()
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
