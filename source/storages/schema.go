package storages

type MusicSource = int64

type Language = int64

const (
	UnknownSource MusicSource = iota
	YouTubeSource
	SpotifySource
)

const (
	LangEnglish Language = iota
	LangFrench
	LangChinese
	LangJapanese
)

// UserProfile table.
type UserProfile struct {
	UserId           string   `db:"user_id" json:"userId"`
	Username         string   `db:"username" json:"username"`
	Language         Language `db:"language" json:"language"`
	RegistrationDate int64    `db:"registration_date" json:"registrationDate"`
	Token            string   `db:"token" json:"token"`
}

// Playlist table.
type Playlist struct {
	UserId       string `db:"user_id" json:"userId"`
	PlaylistId   int64  `db:"playlist_id" json:"playlistId"`
	Title        string `db:"title" json:"title"`
	ModifiedDate int64  `db:"modified_date" json:"modifiedDate"`
	CoverBlob    []byte `db:"cover_blob" json:"coverBlob"`
}

// Music table.
type Music struct {
	MusicId       string      `db:"music_id" json:"musicId"`
	Source        MusicSource `db:"source" json:"source"`
	Title         string      `db:"title" json:"title"`
	LengthSeconds int64       `db:"length_seconds" json:"lengthSeconds"`
}

// PlaylistMusic table.
type PlaylistMusic struct {
	UserId       string `db:"user_id" json:"userId"`
	PlaylistId   int64  `db:"playlist_id" json:"playlistId"`
	MusicId      string `db:"music_id" json:"musicId"`
	Source       int64  `db:"source" json:"source"`
	ModifiedDate int64  `db:"modified_date" json:"modifiedDate"`
}
