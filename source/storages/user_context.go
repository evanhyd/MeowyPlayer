package storages

import (
	"io"
	"log/slog"
	"net/http"
)

type UserContext struct {
	httpClient *http.Client
	config     UserConfig
	storage    Storage
	dispatcher EventsDispatcher
}

func MakeUserContext(httpClient *http.Client, config UserConfig, storage Storage) UserContext {
	return UserContext{
		httpClient: httpClient,
		config:     config,
		storage:    storage,
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

func (u *UserContext) Close() {
	u.httpClient.CloseIdleConnections()

	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			slog.Error("failed to close the storage", "error", err)
		}
	}
}

// DB wrapper calls.
func (u *UserContext) OnInit() {
	u.dispatcher.Dispatch(OnInitEvent, OnInitEventData{})
}

func (u *UserContext) PutUser(userProfile UserProfile) error {
	err := u.storage.PutUser(userProfile)
	if err == nil {
		u.dispatcher.Dispatch(OnPutUserEvent, OnPutUserEventData{UserProfile: userProfile})
	}
	return err
}

func (u *UserContext) DeleteUser() error {
	err := u.storage.DeleteUser()
	if err == nil {
		u.dispatcher.Dispatch(OnDeleteUserEvent, OnDeleteUserEventData{})
	}
	return err
}

func (u *UserContext) GetUser() (UserProfile, error) {
	return u.storage.GetUser()
}

func (u *UserContext) PutPlaylist(playlist Playlist) (Playlist, error) {
	playlist, err := u.storage.PutPlaylist(playlist)
	if err == nil {
		u.dispatcher.Dispatch(OnPutPlaylistEvent, OnPutPlaylistEventData{Playlist: playlist})
	}
	return playlist, err
}

func (u *UserContext) GetPlaylist(playlistID int64) (Playlist, error) {
	return u.storage.GetPlaylist(playlistID)
}

func (u *UserContext) DeletePlaylist(playlistID int64) error {
	err := u.storage.DeletePlaylist(playlistID)
	if err == nil {
		u.dispatcher.Dispatch(OnDeletePlaylistEvent, OnDeletePlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

func (u *UserContext) PutMusic(music Music) error {
	return u.storage.PutMusic(music)
}

func (u *UserContext) GetMusic(musicID string, source MusicSource) (Music, error) {
	return u.storage.GetMusic(musicID, source)
}

func (u *UserContext) DeleteMusic(musicID string, source MusicSource) error {
	return u.storage.DeleteMusic(musicID, source)
}

func (u *UserContext) PutMusicFile(music Music, content io.Reader) error {
	return u.storage.PutMusicFile(music, content)
}

func (u *UserContext) GetMusicFile(music Music) (io.ReadSeekCloser, error) {
	return u.storage.GetMusicFile(music)
}

func (u *UserContext) DeleteMusicFile(music Music) error {
	return u.storage.DeleteMusicFile(music)
}

func (u *UserContext) GetPlaylistsFromUser() ([]Playlist, error) {
	return u.storage.GetPlaylistsFromUser()
}

func (u *UserContext) GetMusicFromPlaylist(playlistID int64) ([]Music, error) {
	return u.storage.GetMusicFromPlaylist(playlistID)
}

func (u *UserContext) PutMusicInPlaylist(playlistID int64, musicID string, source MusicSource) error {
	err := u.storage.PutMusicInPlaylist(playlistID, musicID, source)
	if err == nil {
		u.dispatcher.Dispatch(OnPutMusicInPlaylistEvent, OnPutMusicInPlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

func (u *UserContext) DeleteMusicFromPlaylist(playlistID int64, musicID string, source MusicSource) error {
	err := u.storage.DeleteMusicFromPlaylist(playlistID, musicID, source)
	if err == nil {
		u.dispatcher.Dispatch(OnDeleteMusicFromPlaylistEvent, OnDeleteMusicFromPlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

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
