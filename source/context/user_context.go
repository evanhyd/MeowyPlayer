package context

import "meowyplayer/storage"

const (
	dbPath        = "local.db"
	musicFileBase = "music"
)

type UserContext struct {
	storageManager storage.StorageManager
	dispatcher     EventDispatcher
}

func makeUserContext() UserContext {
	return UserContext{
		storageManager: storage.NewSQLiteStorageManager(dbPath, musicFileBase, storage.User{UserID: 0}),
		dispatcher:     makeEventDispatcher(),
	}
}

func (u *UserContext) AddListener(event Event, listener EventListener) {
	u.dispatcher.addListener(event, listener)
}

func (u *UserContext) LoadPlaylistToMusicPlayer(playlistID int64) error {
	songs, err := u.storageManager.ListMusicInPlaylist(playlistID)
	if err != nil {
		return err
	}

	u.dispatcher.dispatch(OnLoadPlaylistToMusicPlayerEvent, songs)
	return nil
}
