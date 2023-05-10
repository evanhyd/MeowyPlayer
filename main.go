package main

import (
	"log"

	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui"
)

func init() {
}

func main() {
	mainWindow := ui.NewMeowyPlayerWindow()

	meowyPlayerState := player.GetState()
	meowyPlayer := player.GetPlayer()

	if err := player.RefreshAlbumTab(); err != nil {
		log.Fatal(err)
	}
	meowyPlayerState.OnUpdateSeeker().AddCallback(meowyPlayer.SetMusic)

	go meowyPlayer.Launch()
	mainWindow.ShowAndRun()

	player.RemoveUnusedMusic()
}
