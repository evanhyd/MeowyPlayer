package storages

type EventType int64

type EventListener = func(any)

type EventsDispatcher struct {
	listeners map[EventType][]EventListener
}

func makeEventsDispatcher() EventsDispatcher {
	return EventsDispatcher{listeners: make(map[EventType][]EventListener)}
}

func (e *EventsDispatcher) AddListener(event EventType, listener EventListener) {
	e.listeners[event] = append(e.listeners[event], listener)
}

func (e *EventsDispatcher) Dispatch(event EventType, eventData any) {
	for _, listener := range e.listeners[event] {
		listener(eventData)
	}
}

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
