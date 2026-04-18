package context

import "meowyplayer/storages"

type EventType int64
type EventListener = func(EventType, any)

const (
	OnSetStorageEvent EventType = iota
	OnViewPlaylistEvent
	OnReturnBackFromPlaylistEvent
	OnCreatePlaylistEvent
	OnAddMusicToPlaylistEvent
)

type OnSetStorageEventData struct {
	Storage storages.Storage
}

type OnViewPlaylistEventData struct {
	Playlist storages.Playlist
}

type OnReturnBackFromPlaylistEventData struct {
}

type OnCreatePlaylistEventData struct {
	Playlist storages.Playlist
}

type OnAddMusicToPlaylistEventData struct {
	Playlist storages.Playlist
}
