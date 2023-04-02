package main

import (
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui"
)

func main() {
	mainWindow := ui.NewMeowyPlayerWindow()

	meowyPlayerState := player.GetState()
	meowyPlayer := player.GetPlayer()
	meowyPlayerState.OnReadAlbumsFromDisk().NotifyAll(player.ReadAlbumsFromDisk())
	meowyPlayerState.OnSelectMusic().AddObserver(meowyPlayer.PlayerMusicUpdater())
	go meowyPlayer.Launch()

	mainWindow.ShowAndRun()
}
