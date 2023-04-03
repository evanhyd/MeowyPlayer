package main

import (
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui"
)

func main() {
	mainWindow := ui.NewMeowyPlayerWindow()

	meowyPlayerState := player.GetState()
	meowyPlayer := player.GetPlayer()
	meowyPlayerState.OnReadAlbumsFromDiskSubject().NotifyAll(player.ReadAlbumsFromDisk())
	meowyPlayerState.OnSelectMusicSubject().AddCallback(meowyPlayer.SetMusic)
	go meowyPlayer.Launch()

	mainWindow.ShowAndRun()
}
