package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/ui/cwidget"
)

func newController() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	controller := cwidget.NewMusicController()
	manager.GetCurrentPlay().Attach(coverView)
	manager.GetCurrentPlay().Attach(controller)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
