package main

import (
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui"
)

func main() {
	mainWindow := ui.NewMeowyPlayerWindow()
	player.GetPlayerState().UpdateAlbums()
	mainWindow.ShowAndRun()
}
