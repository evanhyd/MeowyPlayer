package mcontext

import "meowyplayer/storages"

type EventType int64
type EventListener = func(any)

const (
	OnInitEvent EventType = iota
	OnPutPlaylistEvent
	OnDeletePlaylistEvent
	OnPutMusicInPlaylistEvent
	OnDeleteMusicFromPlaylistEvent
	OnViewPlaylistPageEvent
	OnViewMusicPageEvent
	OnPlayMusicEvent
	OnPutUserEvent
	OnDeleteUserEvent
)

type OnInitEventData struct {
}

type OnPutPlaylistEventData struct {
	Playlist storages.Playlist
}

type OnDeletePlaylistEventData struct {
	PlaylistId int64
}

type OnPutMusicInPlaylistEventData struct {
	PlaylistId int64
}

type OnDeleteMusicFromPlaylistEventData struct {
	PlaylistId int64
}

type OnViewPlaylistPageEventData struct {
}

type OnViewMusicPageEventData struct {
	Playlist storages.Playlist
}

type OnPlayMusicEventData struct {
	Playlist storages.Playlist
	Music    storages.Music
}

type OnPutUserEventData struct {
	UserProfile storages.UserProfile
}

type OnDeleteUserEventData struct {
}
