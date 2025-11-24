package context

import (
	"meowyplayer/events"
	"meowyplayer/loggers"
	"meowyplayer/storages"
)

type UserContext struct {
	storage    storages.Storage
	dispatcher events.EventsDispatcher
	logger     loggers.Logger
}

func MakeUserContext() UserContext {
	return UserContext{
		dispatcher: events.MakeEventsDispatcher(),
		logger:     loggers.MakeLogger(),
	}
}

func (u *UserContext) AddListener(event events.EventType, listener events.EventListener) {
	u.dispatcher.AddListener(event, listener)
}

func (u *UserContext) Close() {
	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			u.logger.Log.Error("failed to close the storage", "error", err)
		}
	}

	if err := u.logger.Close(); err != nil {
		u.logger.Log.Error("failed to close the logger", "error", err)
	}
}

func (u *UserContext) SetStorage(storage storages.Storage) {
	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			u.logger.Log.Error("failed to close the storage", "error", err)
		}
	}
	u.storage = storage
	u.dispatcher.Dispatch(events.StorageSetEvent, events.StorageSetEventData{u.storage})
}

func (u *UserContext) SetPlaylist(playlistID int64, index int) error {
	musicList, err := u.storage.ListMusicInPlaylist(playlistID)
	if err != nil {
		return err
	}

	u.dispatcher.Dispatch(events.PlaylistSetEvent, events.PlaylistSetEventData{musicList, index})
	return nil
}
