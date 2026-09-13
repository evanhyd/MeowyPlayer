package storages

import (
	"log/slog"
	"net/http"
)

type UserContext struct {
	Storage // Anonymous embedding: inherits all storage methods automatically

	httpClient *http.Client
	config     UserConfig
	dispatcher EventsDispatcher
}

func MakeUserContext(httpClient *http.Client, config UserConfig, storage Storage) UserContext {
	return UserContext{
		Storage:    storage,
		httpClient: httpClient,
		config:     config,
		dispatcher: makeEventsDispatcher(),
	}
}

func (u *UserContext) AddListener(event EventType, listener EventListener) {
	u.dispatcher.AddListener(event, listener)
}

func (u *UserContext) HttpClient() *http.Client {
	return u.httpClient
}

func (u *UserContext) Config() *UserConfig {
	return &u.config
}

// ---------------- Overridden Methods (Custom Logic / Event Dispatching) ----------------

func (u *UserContext) Close() {
	u.httpClient.CloseIdleConnections()

	if u.Storage != nil {
		if err := u.Storage.Close(); err != nil {
			slog.Error("failed to close the storage", "error", err)
		}
	}
}

func (u *UserContext) OnInit() {
	u.dispatcher.Dispatch(OnInitEvent, OnInitEventData{})
}

func (u *UserContext) PutUser(userProfile UserProfile) error {
	err := u.Storage.PutUser(userProfile)
	if err == nil {
		u.dispatcher.Dispatch(OnPutUserEvent, OnPutUserEventData{UserProfile: userProfile})
	}
	return err
}

func (u *UserContext) DeleteUser() error {
	err := u.Storage.DeleteUser()
	if err == nil {
		u.dispatcher.Dispatch(OnDeleteUserEvent, OnDeleteUserEventData{})
	}
	return err
}

func (u *UserContext) PutPlaylist(playlist Playlist) (Playlist, error) {
	playlist, err := u.Storage.PutPlaylist(playlist)
	if err == nil {
		u.dispatcher.Dispatch(OnPutPlaylistEvent, OnPutPlaylistEventData{Playlist: playlist})
	}
	return playlist, err
}

func (u *UserContext) DeletePlaylist(playlistID int64) error {
	err := u.Storage.DeletePlaylist(playlistID)
	if err == nil {
		u.dispatcher.Dispatch(OnDeletePlaylistEvent, OnDeletePlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

func (u *UserContext) PutPlaylistMusic(playlistMusic PlaylistMusic) error {
	err := u.Storage.PutPlaylistMusic(playlistMusic)
	if err == nil {
		u.dispatcher.Dispatch(OnPutMusicInPlaylistEvent, OnPutMusicInPlaylistEventData{PlaylistId: playlistMusic.PlaylistId})
	}
	return err
}

func (u *UserContext) DeletePlaylistMusic(playlistID int64, musicID string, source MusicSource) error {
	err := u.Storage.DeletePlaylistMusic(playlistID, musicID, source)
	if err == nil {
		u.dispatcher.Dispatch(OnDeleteMusicFromPlaylistEvent, OnDeleteMusicFromPlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

// ---------------- UI Event Dispatching ----------------

func (u *UserContext) ViewPlaylistPage() {
	u.dispatcher.Dispatch(OnViewPlaylistPageEvent, OnViewPlaylistPageEventData{})
}

func (u *UserContext) ViewMusicPage(playlist Playlist) {
	u.dispatcher.Dispatch(OnViewMusicPageEvent, OnViewMusicPageEventData{Playlist: playlist})
}

func (u *UserContext) PlayMusic(playlist Playlist, music Music) error {
	u.dispatcher.Dispatch(OnPlayMusicEvent, OnPlayMusicEventData{Playlist: playlist, Music: music})
	return nil
}
