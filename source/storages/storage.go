package storages

import (
	"io"
)

// All interface implementation must be idempotent.
// Put creates or updates the data.
// Delete deletes the data or do nothing if the data doesn't exist.

type PlaylistStorer interface {
	PutPlaylist(session UserSession, playlist Playlist) (Playlist, error)
	GetPlaylist(session UserSession, playlistID int64) (Playlist, error)
	DeletePlaylist(session UserSession, playlistID int64) error
}

type MusicStorer interface {
	PutMusic(session UserSession, music Music) error
	GetMusic(session UserSession, musicID string, source MusicSource) (Music, error)
	DeleteMusic(session UserSession, musicID string, source MusicSource) error
}

type FileStorer interface {
	PutMusicFile(session UserSession, music Music, content io.Reader) error
	GetMusicFile(session UserSession, music Music) (io.ReadCloser, error)
	DeleteMusicFile(session UserSession, music Music) error
}

type PlaylistManager interface {
	GetPlaylistsFromUser(session UserSession) ([]Playlist, error)
	GetMusicFromPlaylist(session UserSession, playlistID int64) ([]Music, error)
	PutMusicInPlaylist(session UserSession, playlistID int64, musicID string, source MusicSource) error
	DeleteMusicFromPlaylist(session UserSession, playlistID int64, musicID string, source MusicSource) error
}

type Storage interface {
	PlaylistStorer
	MusicStorer
	FileStorer
	PlaylistManager
	io.Closer
}
