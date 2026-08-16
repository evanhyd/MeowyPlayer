package storages

type MusicSource int64

const (
	UnknownSource MusicSource = iota
	YouTubeSource
	SpotifySource
)

type Language int64

const (
	LangEnglish Language = iota
	LangFrench
	LangChinese
	LangJapanese
)

type UserProfile struct {
	UserId           string   `db:"user_id"`
	Username         string   `db:"username"`
	Language         Language `db:"language"`
	RegistrationDate int64    `db:"registration_date"`
	Token            string   `db:"token"`
}

type Playlist struct {
	UserId       string `db:"user_id"`
	PlaylistId   int64  `db:"playlist_id"`
	Deleted      bool   `db:"deleted"`
	Title        string `db:"title"`
	ModifiedDate int64  `db:"modified_date"`
	CoverBlob    []byte `db:"cover_blob"`
}

type Music struct {
	MusicId       string      `db:"music_id"`
	Source        MusicSource `db:"source"`
	Title         string      `db:"title"`
	LengthSeconds int64       `db:"length_seconds"`
}

type PlaylistMusic struct {
	UserId     string `db:"user_id"`
	PlaylistId int64  `db:"playlist_id"`
	MusicId    string `db:"music_id"`
	Source     int64  `db:"source"`
	AddedAt    int64  `db:"added_at"`
}
