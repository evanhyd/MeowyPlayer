package ui

import (
	"fmt"

	"meowyplayer.com/core/client"
)

func showAddLocalAlbumDialog() {
	showErrorIfAny(client.AddRandomAlbum())
}

func showAddOnlineAlbumDialog() {
	fmt.Println("not completed")
}
