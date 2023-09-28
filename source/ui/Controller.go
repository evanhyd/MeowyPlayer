package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/source/utility"
)

func newController() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	musicPlayer := client.NewMusicPlayer()
	controller := cwidget.NewPlayerMenu()

	client.GetCurrentPlayList().Attach(utility.MakeCallback(func(p *player.PlayList) {
		coverView.SetAlbum(p.Album())
		coverView.OnTapped = func(*fyne.PointEvent) { client.GetCurrentAlbum().Set(p.Album()) }
	}))
	client.GetCurrentPlayList().Attach(musicPlayer)
	go musicPlayer.Start(controller)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
