package main

import (
	"meowyplayer/context"
	"meowyplayer/events"
	"meowyplayer/players"
	"meowyplayer/storages"
)

func main() {
	const (
		dbPath        = "local.db"
		musicFileBase = "music"
	)

	player := players.MakeBeepPlayer()
	userContext := context.MakeUserContext()
	userContext.AddListener(events.StorageSetEvent, player.HandleStorageSetEvent)
	userContext.AddListener(events.PlaylistSetEvent, player.HandlePlaylistSetEvent)
	userContext.SetStorage(storages.NewSQLiteStorage(dbPath, musicFileBase, storages.User{UserID: 0}))

}
