package storages

import (
	"strings"
)

type MusicSource int64

const (
	UnknownSource MusicSource = iota
	YouTubeSource
	SpotifySource
)

// Users table
type User struct {
	UserId         int64
	Name           string
	HashedPassword string
	Salt           string
}

// UserSessions table
type UserSession struct {
	UserId    int64
	Token     string
	CreatedAt int64 // Unix nano
}

// Playlist table
type Playlist struct {
	UserId       int64
	PlaylistId   int64
	Title        string
	ModifiedDate int64 // Unix nano
	CoverBlob    []byte
}

// Music table
type Music struct {
	MusicId       string
	Source        MusicSource
	Title         string
	LengthSeconds int64
}

// PlaylistMusic table
type PlaylistMusic struct {
	UserId     int64
	PlaylistId int64
	MusicId    string
	Source     MusicSource
	AddedAt    int64 // Unix nano
}

func ComparePlaylist(l Playlist, r Playlist) int {
	if l.UserId != r.UserId {
		return int(l.UserId) - int(r.UserId)
	}
	return int(l.PlaylistId) - int(r.PlaylistId)
}

func CompareMusic(l Music, r Music) int {
	if l.MusicId != r.MusicId {
		return strings.Compare(l.MusicId, r.MusicId)
	}
	return int(l.Source) - int(r.Source)
}
