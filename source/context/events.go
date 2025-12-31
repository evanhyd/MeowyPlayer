package context

import "meowyplayer/storages"

type EventType int64
type EventListener = func(EventType, any)

const (
	StorageSetEvent EventType = iota
	PlaylistSetEvent
)

type StorageSetEventData struct {
	Storage storages.Storage
}

type PlaylistSetEventData struct {
	MusicList []storages.Music
	Index     int
}
