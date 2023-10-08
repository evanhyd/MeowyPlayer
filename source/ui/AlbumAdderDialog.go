package ui

import (
	"fmt"

	"meowyplayer.com/source/client"
)

func showAddLocalAlbumDialog() {
	showErrorIfAny(client.AddAlbum())
}

func showAddOnlineAlbumDialog() {
	fmt.Println("not completed")
}
