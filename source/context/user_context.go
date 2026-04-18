package context

import (
	"log/slog"
	"meowyplayer/storages"
)

type UserContext struct {
	storage    storages.Storage
	dispatcher EventsDispatcher
}

func MakeUserContext() UserContext {
	return UserContext{dispatcher: makeEventsDispatcher()}
}

func (u *UserContext) AddListener(event EventType, listener EventListener) {
	u.dispatcher.AddListener(event, listener)
}

func (u *UserContext) Close() {
	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			slog.Error("failed to close the storage", "error", err)
		}
	}
}

func (u *UserContext) Storage() storages.Storage {
	return u.storage
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

func (u *UserContext) SetPlaylist(playlistID int64, index int) error {
	musicList, err := u.storage.GetAllMusicFromPlaylist(playlistID)
	if err != nil {
		return err
	}

	u.dispatcher.Dispatch(OnSetPlaylistEvent, OnSetPlaylistEventData{MusicList: musicList, Index: index})
	return nil
}

func (u *UserContext) CreatePlaylist(title string, coverBlob []byte) error {
	playlist, err := u.storage.CreatePlaylist(title, coverBlob)
	if err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnCreatePlaylistEvent, OnCreatePlaylistEventData{Playlist: playlist})
	return nil
}

func (u *UserContext) AddMusicToPlaylist(playlist storages.Playlist, musicID string, source storages.MusicSource) error {
	err := u.storage.AddMusicToPlaylist(playlist.PlaylistId, musicID, source)
	if err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnAddMusicToPlaylistEvent, OnAddMusicToPlaylistEventData{Playlist: playlist})
	return nil
}
