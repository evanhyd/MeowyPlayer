package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/logger"
	"meowyplayer.com/utility/pattern"
)

func NewMainWindow() fyne.Window {
	//create window
	window := fyne.CurrentApp().NewWindow("MeowyPlayer")
	window.SetCloseIntercept(window.Hide)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(770.0, 650.0))

	//create item tabs
	albumTab := newAlbumTab()
	musicTab := newMusicTab()
	accountTab := newClientTab()
	tabs := container.NewAppTabs(albumTab, musicTab, accountTab)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.DisableItem(musicTab)

	client.Manager().AddAlbumListener(pattern.MakeCallback(func(resource.Album) {
		tabs.EnableItem(musicTab)
		tabs.Select(musicTab)
	}))

	window.SetContent(container.NewBorder(nil, newController(), nil, nil, tabs))
	return window
}

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func showErrorIfAny(err error) {
	if err != nil {
		logger.Error(err, 1)
		dialog.ShowError(err, getWindow())
	}
}
