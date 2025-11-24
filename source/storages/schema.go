package storages

type MusicSource int64

const (
	UnknownSource MusicSource = iota
	YouTubeSource
	SpotifySource
)

// Users table
type User struct {
	UserID         int64
	Name           string
	HashedPassword string
	Salt           string
}

// UserSessions table
type UserSession struct {
	UserID    int64
	Token     string
	CreatedAt int64 // Unix nano
}

// Playlist table
type Playlist struct {
	PlaylistID int64
	Title      string
	CoverBlob  []byte
}

// Music table
type Music struct {
	MusicID       string
	Source        MusicSource
	Title         string
	LengthSeconds int64
}

// PlaylistMusic table
type PlaylistMusic struct {
	PlaylistID int64
	MusicID    string
	Source     MusicSource
}
