package context

import (
	"io"
	"log/slog"
	"meowyplayer/storages"
)

type UserContext struct {
	session    storages.UserSession
	storage    storages.Storage
	dispatcher EventsDispatcher
}

func MakeUserContext() UserContext {
	return UserContext{dispatcher: makeEventsDispatcher()}
}

func (u *UserContext) AddListener(event EventType, listener EventListener) {
	u.dispatcher.AddListener(event, listener)
}

func (u *UserContext) SetStorage(storage storages.Storage) {
	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			slog.Error("failed to close the storage", "error", err)
		}
	}
	u.storage = storage
	u.dispatcher.Dispatch(OnSetStorageEvent, OnSetStorageEventData{Storage: u.storage})
}

func (u *UserContext) Close() {
	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			slog.Error("failed to close the storage", "error", err)
		}
	}
}

// DB wrapper calls.
func (u *UserContext) PutPlaylist(playlist storages.Playlist) (storages.Playlist, error) {
	playlist, err := u.storage.PutPlaylist(u.session, playlist)
	if err == nil {
		u.dispatcher.Dispatch(OnPutPlaylistEvent, OnPutPlaylistEventData{Playlist: playlist})
	}
	return playlist, err
}

func (u *UserContext) GetPlaylist(playlistID int64) (storages.Playlist, error) {
	return u.storage.GetPlaylist(u.session, playlistID)
}

func (u *UserContext) DeletePlaylist(playlistID int64) error {
	err := u.storage.DeletePlaylist(u.session, playlistID)
	if err == nil {
		u.dispatcher.Dispatch(OnDeletePlaylistEvent, OnDeletePlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

func (u *UserContext) PutMusic(music storages.Music) error {
	return u.storage.PutMusic(u.session, music)
}

func (u *UserContext) GetMusic(musicID string, source storages.MusicSource) (storages.Music, error) {
	return u.storage.GetMusic(u.session, musicID, source)
}

func (u *UserContext) DeleteMusic(musicID string, source storages.MusicSource) error {
	return u.storage.DeleteMusic(u.session, musicID, source)
}

func (u *UserContext) PutMusicFile(music storages.Music, content io.Reader) error {
	return u.storage.PutMusicFile(u.session, music, content)
}

func (u *UserContext) GetMusicFile(music storages.Music) (io.ReadCloser, error) {
	return u.storage.GetMusicFile(u.session, music)
}

func (u *UserContext) DeleteMusicFile(music storages.Music) error {
	return u.storage.DeleteMusicFile(u.session, music)
}

func (u *UserContext) GetPlaylistsFromUser() ([]storages.Playlist, error) {
	return u.storage.GetPlaylistsFromUser(u.session)
}

func (u *UserContext) GetMusicFromPlaylist(playlistID int64) ([]storages.Music, error) {
	return u.storage.GetMusicFromPlaylist(u.session, playlistID)
}

func (u *UserContext) PutMusicInPlaylist(playlistID int64, musicID string, source storages.MusicSource) error {
	err := u.storage.PutMusicInPlaylist(u.session, playlistID, musicID, source)
	if err == nil {
		u.dispatcher.Dispatch(OnPutMusicInPlaylistEvent, OnPutMusicInPlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

func (u *UserContext) DeleteMusicFromPlaylist(playlistID int64, musicID string, source storages.MusicSource) error {
	err := u.storage.DeleteMusicFromPlaylist(u.session, playlistID, musicID, source)
	if err == nil {
		u.dispatcher.Dispatch(OnDeleteMusicFromPlaylistEvent, OnDeleteMusicFromPlaylistEventData{PlaylistId: playlistID})
	}
	return err
}

func (u *UserContext) ViewPlaylistPage() {
	u.dispatcher.Dispatch(OnViewPlaylistPageEvent, OnViewPlaylistPageEventData{})
}

func (u *UserContext) ViewMusicPage(playlist storages.Playlist) {
	u.dispatcher.Dispatch(OnViewMusicPageEvent, OnViewMusicPageEventData{Playlist: playlist})
}

func (u *UserContext) PlayMusic(playlist storages.Playlist, music storages.Music) error {
	u.dispatcher.Dispatch(OnPlayMusicEvent, OnPlayMusicEventData{Playlist: playlist, Music: music})
	return nil
}
