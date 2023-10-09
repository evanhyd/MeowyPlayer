package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/utility/pattern"
)

func NewMainWindow() fyne.Window {
	const windowTitle = "MeowyPlayer"
	windowSize := fyne.NewSize(770.0, 650.0)

	//create window
	fyne.SetCurrentApp(app.NewWithID(windowTitle))
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	window := newWindow(windowTitle, windowSize)

	//create system tray
	if desktop, ok := fyne.CurrentApp().(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
		desktop.SetSystemTrayIcon(resource.WindowIcon())
	}

	//create item tabs
	albumTab := newAlbumTab()
	musicTab := newMusicTab()
	tabs := container.NewAppTabs(albumTab, musicTab)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.DisableItem(musicTab)
	client.GetAlbumData().Attach(pattern.MakeCallback(func(*resource.Album) {
		tabs.EnableItem(musicTab)
		tabs.Select(musicTab)
	}))

	window.SetContent(container.NewBorder(nil, newController(), nil, nil, tabs))
	return window
}

func newWindow(title string, size fyne.Size) fyne.Window {
	window := fyne.CurrentApp().NewWindow(title)
	window.SetMaster()
	window.SetIcon(resource.WindowIcon())
	window.SetCloseIntercept(window.Hide)
	window.Resize(size)
	window.CenterOnScreen()
	return window
}

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func showErrorIfAny(err error) {
	if err != nil {
		dialog.ShowError(err, getWindow())
	}
}
