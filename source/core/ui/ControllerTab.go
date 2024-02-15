package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/player"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newControllerTab() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	controller := cwidget.NewMusicController()
	musicPlayer := player.NewMP3Player()
	controller.Bind(musicPlayer)
	go musicPlayer.Start(controller)

	// set album preview when loading a new playlist
	client.Manager().AddPlayListListener(pattern.MakeCallback(func(player.PlayList) {
		focused := client.Manager().FocusedAlbum()
		coverView.SetAlbum(&focused)
		coverView.OnTapped = func(*fyne.PointEvent) {
			showErrorIfAny(client.Manager().SetFocusedAlbum(focused))
			//TODO: fix reference when renaming/deleting the album
		}
	}))

	// load playlist to the mp3 player
	client.Manager().AddPlayListListener(musicPlayer)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
