package storages

import (
	"io"
)

type UserStorer interface {
	PutUser(userProfile UserProfile) error
	GetUser() (UserProfile, error)
	DeleteUser() error
}

// UserId will always get overriden by the current user from GetUser().
type PlaylistStorer interface {
	PutPlaylist(playlist Playlist) (Playlist, error)
	GetPlaylist(playlistId int64) (Playlist, error)
	GetPlaylists() ([]Playlist, error)
	DeletePlaylist(playlistId int64) error

	PutPlaylistMusic(playlistMusic PlaylistMusic) error
	GetAllPlaylistMusic(playlistId int64) ([]PlaylistMusic, error)
	DeletePlaylistMusic(playlistMusic PlaylistMusic) error
}

type MusicStorer interface {
	PutMusic(music Music) error
	GetMusic(musicId string, source MusicSource) (Music, error)
	GetAllMusic(playlistId int64) ([]Music, error)
	DeleteMusic(musicId string, source MusicSource) error
}

type FileStorer interface {
	PutMusicFile(music Music, content io.Reader) error
	GetMusicFile(music Music) (io.ReadSeekCloser, error)
	DeleteMusicFile(music Music) error
}

type Storage interface {
	UserStorer
	PlaylistStorer
	MusicStorer
	FileStorer
	io.Closer
}
