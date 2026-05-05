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

/*
Helpers to update the storage and the UI.
*/

func (u *UserContext) EnterPlaylist(playlist storages.Playlist) {
	u.dispatcher.Dispatch(OnViewPlaylistEvent, OnEnterPlaylist{Playlist: playlist})
}

func (u *UserContext) ExitPlaylist() {
	u.dispatcher.Dispatch(OnReturnBackFromPlaylistEvent, OnExitPlaylist{})
}

func (u *UserContext) CreatePlaylist(title string, coverBlob []byte) error {
	playlist, err := u.storage.CreatePlaylist(title, coverBlob)
	if err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnCreatePlaylistEvent, OnCreatePlaylistEventData{Playlist: playlist})
	return nil
}

func (u *UserContext) AddMusic(playlist storages.Playlist, music storages.Music) error {
	err := u.storage.CreateMusic(music)
	if err != nil {
		return err
	}
	err = u.storage.AddMusicToPlaylist(playlist.PlaylistId, music.MusicId, music.Source)
	if err != nil {
		return err
	}

	u.dispatcher.Dispatch(OnAddMusicToPlaylistEvent, OnAddMusicToPlaylistEventData{Playlist: playlist})
	return nil
}

func (u *UserContext) UpdatePlaylist(playlist storages.Playlist) error {
	if err := u.storage.UpdatePlaylist(playlist); err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnUpdatePlaylistEvent, OnUpdatePlaylistEventData{Playlist: playlist})
	return nil
}

func (u *UserContext) DeletePlaylist(playlist storages.Playlist) error {
	if err := u.storage.DeletePlaylist(playlist.PlaylistId); err != nil {
		return err
	}
	u.dispatcher.Dispatch(OnDeletePlaylistEvent, OnDeletePlaylistEventData{})
	return nil
}

func (u *UserContext) PlayPlaylist(playlist storages.Playlist, selectedMusic storages.Music) error {
	u.dispatcher.Dispatch(OnPlayPlaylistEvent, OnPlayPlaylistEventData{Playlist: playlist, SelectedMusic: selectedMusic})
	return nil
}
