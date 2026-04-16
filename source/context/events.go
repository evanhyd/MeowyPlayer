package context

import "meowyplayer/storages"

type EventType int64
type EventListener = func(EventType, any)

const (
	OnSetStorageEvent EventType = iota
	OnSetPlaylistEvent
	OnCreatePlaylistEvent
	OnAddMusicToPlaylistEvent
)

type OnSetStorageEventData struct {
	Storage storages.Storage
}

type OnSetPlaylistEventData struct {
	MusicList []storages.Music
	Index     int
}

type OnCreatePlaylistEventData struct {
	PlaylistID int64
}

type OnAddMusicToPlaylistEventData struct {
	PlaylistID int64
}
