package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newController() fyne.CanvasObject {
	coverView := cwidget.NewCoverView(fyne.NewSize(128.0, 128.0))
	controller := cwidget.NewMediaMenu()
	musicPlayer := client.NewMusicPlayer()
	controller.Bind(musicPlayer)
	go musicPlayer.Start(controller)

	client.GetInstance().AddPlayListListener(pattern.MakeCallback(func(p resource.PlayList) {
		album := client.GetInstance().GetAlbum()
		coverView.SetAlbum(&album)
		coverView.OnTapped = func(*fyne.PointEvent) { client.GetInstance().SetAlbum(album) }
	}))
	client.GetInstance().AddPlayListListener(musicPlayer)

	return container.NewBorder(nil, nil, coverView, nil, controller)
}
