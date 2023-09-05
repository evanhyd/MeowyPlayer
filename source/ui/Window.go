package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"meowyplayer.com/source/resource/texture"
)

func NewMainWindow() fyne.Window {
	const (
		iconName    = "icon.ico"
		appID       = "MeowyPlayer"
		windowTitle = "MeowyPlayer"
	)
	windowSize := fyne.NewSize(460.0, 650.0)
	icon := texture.Get(iconName)

	fyne.SetCurrentApp(app.NewWithID(appID))
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
	tabs := container.NewAppTabs(newAlbumTab(), newMusicTab())
	tabs.SetTabLocation(container.TabLocationLeading)
	window.SetContent(container.NewBorder(nil, nil, nil, nil, tabs))
	return window
}

func showErrorIfAny(err error) {
	if err != nil {
		log.Printf("show error: %v\n", err)
		dialog.ShowError(err, getMainWindow())
	}
}

func getMainWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}
