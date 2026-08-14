package storages

import (
	"io"
)

type UserStorer interface {
	PutUser(userProfile UserProfile) error
	GetUser() (UserProfile, error)
	DeleteUser() error
}

type PlaylistStorer interface {
	PutPlaylist(playlist Playlist) (Playlist, error)
	GetPlaylist(playlistId int64) (Playlist, error)
	DeletePlaylist(playlistId int64) error
}

type MusicStorer interface {
	PutMusic(music Music) error
	GetMusic(musicId string, source MusicSource) (Music, error)
	DeleteMusic(musicId string, source MusicSource) error
}

type FileStorer interface {
	PutMusicFile(music Music, content io.Reader) error
	GetMusicFile(music Music) (io.ReadCloser, error)
	DeleteMusicFile(music Music) error
}

type PlaylistManager interface {
	GetPlaylistsFromUser() ([]Playlist, error)
	GetMusicFromPlaylist(playlistId int64) ([]Music, error)
	PutMusicInPlaylist(playlistId int64, musicId string, source MusicSource) error
	DeleteMusicFromPlaylist(playlistId int64, musicId string, source MusicSource) error
}

type Storage interface {
	UserStorer
	PlaylistStorer
	MusicStorer
	FileStorer
	PlaylistManager
	io.Closer
}
