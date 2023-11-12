package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/player"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newController() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	controller := cwidget.NewMediaMenu()
	musicPlayer := player.NewMusicPlayer()
	controller.Bind(musicPlayer)
	go musicPlayer.Start(controller)

	client.Manager().AddPlayListListener(pattern.MakeCallback(func(p player.PlayList) {
		album := client.Manager().Album()
		coverView.SetAlbum(&album)
		coverView.OnTapped = func(*fyne.PointEvent) { client.Manager().SetAlbum(album) }
	}))
	client.Manager().AddPlayListListener(musicPlayer)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
