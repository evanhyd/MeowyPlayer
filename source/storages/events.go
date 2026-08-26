package storages

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
	Playlist Playlist
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
	Playlist Playlist
}

type OnPlayMusicEventData struct {
	Playlist Playlist
	Music    Music
}

type OnPutUserEventData struct {
	UserProfile UserProfile
}

type OnDeleteUserEventData struct {
}
