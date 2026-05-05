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
	OnUpdatePlaylistEvent
	OnDeletePlaylistEvent
	OnPlayPlaylistEvent
	OnDeleteMusicEvent
)

type OnSetStorageEventData struct {
	Storage storages.Storage
}

type OnEnterPlaylist struct {
	Playlist storages.Playlist
}

type OnExitPlaylist struct {
}

type OnCreatePlaylistEventData struct {
	Playlist storages.Playlist
}

type OnAddMusicToPlaylistEventData struct {
	Playlist storages.Playlist
}

type OnUpdatePlaylistEventData struct {
	Playlist storages.Playlist
}

type OnDeletePlaylistEventData struct {
}

type OnPlayPlaylistEventData struct {
	Playlist      storages.Playlist
	SelectedMusic storages.Music
}

type OnDeleteMusicEventData struct {
	Playlist storages.Playlist
}
