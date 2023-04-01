package main

import (
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui"
)

func main() {
	mainWindow := ui.NewMeowyPlayerWindow()
	state := player.GetState()
	state.OnReadAlbums().NotifyAll(player.ReadAlbumsFromDirectory())
	mainWindow.ShowAndRun()
}
