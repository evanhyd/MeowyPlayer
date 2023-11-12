package ui

import (
	"meowyplayer.com/core/client"
)

func showAddLocalAlbumDialog() {
	showErrorIfAny(client.Manager().AddRandomAlbum())
}

func showAddOnlineAlbumDialog() {
	//not implemented
}
