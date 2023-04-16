package main

import (
	"log"

	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui"
)

func main() {
	mainWindow := ui.NewMeowyPlayerWindow()

	meowyPlayerState := player.GetState()
	meowyPlayer := player.GetPlayer()

	albums, err := player.ReadAlbumsFromDisk()
	if err != nil {
		log.Fatal(err)
	}
	meowyPlayerState.OnReadAlbumsFromDiskSubject().NotifyAll(albums)
	meowyPlayerState.OnSelectMusicSubject().AddCallback(meowyPlayer.SetMusic)

	go meowyPlayer.Launch()
	mainWindow.ShowAndRun()
}
