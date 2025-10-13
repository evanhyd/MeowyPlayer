package storage

// StorageManager defines operations for managing playlists, music, and their relationships.
// Assumes user/session context is already encapsulated in the implementation.
type StorageManager interface {
	// --- Playlist ---
	CreatePlaylist(p Playlist) (Playlist, error)
	UpdatePlaylist(p Playlist) error
	DeletePlaylist(playlistID int64) error
	GetPlaylist(playlistID int64) (Playlist, error)
	ListPlaylists() ([]Playlist, error)

	// --- Music ---
	CreateMusic(m Music) error
	UpdateMusic(m Music) error
	DeleteMusic(musicID string, source MusicSource) error
	GetMusic(musicID string, source MusicSource) (Music, error)
	ListMusic() ([]Music, error)

	// --- Playlist - Music association ---
	AddMusicToPlaylist(playlistID int64, musicID string, source MusicSource) error
	RemoveMusicFromPlaylist(playlistID int64, musicID string, source MusicSource) error
	ListMusicInPlaylist(playlistID int64) ([]Music, error)

	Close() error
}
