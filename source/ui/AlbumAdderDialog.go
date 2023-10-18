package ui

import (
	"fmt"

	"meowyplayer.com/source/client"
)

func showAddLocalAlbumDialog() {
	showErrorIfAny(client.AddRandomAlbum())
}

func showAddOnlineAlbumDialog() {
	fmt.Println("not completed")
}
