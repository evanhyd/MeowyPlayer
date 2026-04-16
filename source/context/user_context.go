package context

import (
	"log/slog"
	"meowyplayer/storages"

	"fyne.io/fyne/v2"
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
	musicList, err := u.storage.ListMusicInPlaylist(playlistID)
	if err != nil {
		return err
	}

	u.dispatcher.Dispatch(OnSetPlaylistEvent, OnSetPlaylistEventData{MusicList: musicList, Index: index})
	return nil
}

func (u *UserContext) CreatePlaylist(title string, cover fyne.Resource) error {
	playlist, err := u.storage.CreatePlaylist(storages.Playlist{PlaylistID: 0, Title: title, CoverBlob: cover.Content()})
	if err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnCreatePlaylistEvent, OnCreatePlaylistEventData{PlaylistID: playlist.PlaylistID})
	return nil
}

func (u *UserContext) AddMusicToPlaylist(playlistID int64, musicID string, source storages.MusicSource) error {
	err := u.storage.AddMusicToPlaylist(playlistID, musicID, source)
	if err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnAddMusicToPlaylistEvent, OnAddMusicToPlaylistEventData{PlaylistID: playlistID})
	return nil
}
