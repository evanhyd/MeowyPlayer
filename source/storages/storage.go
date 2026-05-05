package storages

import "io"

// Storage defines operations for managing playlists, music, and their relationships.
// Assumes user/session context is already encapsulated in the implementation.
type Storage interface {
	// Playlist
	CreatePlaylist(title string, coverBlob []byte) (Playlist, error)
	UpdatePlaylist(playlist Playlist) error
	DeletePlaylist(playlistID int64) error
	GetPlaylist(playlistID int64) (Playlist, error)
	GetAllSortedPlaylists() ([]Playlist, error) // By date

	// Music
	CreateMusic(music Music) error
	UpdateMusic(music Music) error
	DeleteMusic(musicID string, source MusicSource) error
	GetMusic(musicID string, source MusicSource) (Music, error)
	GetAllMusic() ([]Music, error)

	// Playlist - Music association
	AddMusicToPlaylist(playlistID int64, musicID string, source MusicSource) error
	RemoveMusicFromPlaylist(playlistID int64, musicID string, source MusicSource) error
	GetAllSortedMusicFromPlaylist(playlistID int64) ([]Music, error) // By date

	// Filesystem manipulation. Does NOT affect the DB table.
	CreateOrUpdateMusicFile(music Music, content io.Reader) error
	RemoveMusicFile(music Music) error
	GetMusicFile(music Music) (io.ReadCloser, error)

	Close() error
}
