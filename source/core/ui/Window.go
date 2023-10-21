package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/pattern"
)

func NewMainWindow() fyne.Window {
	window := newWindow("MeowyPlayer", fyne.NewSize(770.0, 650.0))

	//create item tabs
	albumTab := newAlbumTab()
	musicTab := newMusicTab()
	tabs := container.NewAppTabs(albumTab, musicTab)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.DisableItem(musicTab)

	client.GetInstance().AddAlbumListener(pattern.MakeCallback(func(resource.Album) {
		tabs.EnableItem(musicTab)
		tabs.Select(musicTab)
	}))

	window.SetContent(container.NewBorder(nil, newController(), nil, nil, tabs))
	return window
}

func newWindow(title string, size fyne.Size) fyne.Window {
	window := fyne.CurrentApp().NewWindow(title)
	window.SetMaster()
	window.SetIcon(resource.WindowIcon)
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
		log.Println(err)
	}
}
