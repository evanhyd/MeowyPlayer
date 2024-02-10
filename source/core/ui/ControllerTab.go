package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/player"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newControllerTab() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	controller := cwidget.NewMusicController()
	musicPlayer := player.NewMP3Player()
	controller.Bind(musicPlayer)
	go musicPlayer.Start(controller)

	client.Manager().AddAlbumListener(pattern.MakeCallback(func(album resource.Album) {
		coverView.SetAlbum(&album)
		coverView.OnTapped = func(*fyne.PointEvent) { client.Manager().SetAlbum(album) }
	}))

	client.Manager().AddAlbumPlayedListener(&musicPlayer)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
