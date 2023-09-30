package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newController() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	musicPlayer := client.NewMusicPlayer()
	controller := cwidget.NewPlayerMenu()
	go musicPlayer.Start(controller)

	client.GetPlayListData().Attach(pattern.MakeCallback(func(p *player.PlayList) {
		coverView.SetAlbum(p.Album())
		coverView.OnTapped = func(*fyne.PointEvent) { client.GetAlbumData().Set(p.Album()) }
	}))
	client.GetPlayListData().Attach(musicPlayer)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
