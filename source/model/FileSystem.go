package model

import (
	"encoding/json"
	"io"
	"time"

	"fyne.io/fyne/v2"
)

// random UUID
// watch out path injection
type AlbumKey string

func (k AlbumKey) IsEmpty() bool {
	return len(k) == 0
}

type Album struct {
	key   AlbumKey
	date  time.Time
	title string
	music []Music
	cover []byte
}

type albumJson struct {
	Key   AlbumKey `json:"key"`
	Date  int64    `json:"date"` //epoch seconds
	Title string   `json:"title"`
	Music []Music  `json:"music"`
	Cover []byte   `json:"cover"`
}

func (a *Album) Key() AlbumKey {
	return a.key
}

func (a *Album) Date() time.Time {
	return a.date
}

func (a Album) Title() string {
	return a.title
}

func (a *Album) Music() []Music {
	return a.music
}

func (a *Album) Cover() fyne.Resource {
	return fyne.NewStaticResource("", a.cover)
}

func (a *Album) Count() int {
	return len(a.music)
}

func (a *Album) MarshalJSON() ([]byte, error) {
	return json.Marshal(albumJson{a.key, a.date.Unix(), a.title, a.music, a.cover})
}

func (a *Album) UnmarshalJSON(data []byte) error {
	var jsonData albumJson
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	*a = Album{jsonData.Key, time.Unix(jsonData.Date, 0), jsonData.Title, jsonData.Music, jsonData.Cover}
	return nil
}

// unique hash that maps to file system's file name.
// watch out path injection
type MusicKey string

func (k MusicKey) IsEmpty() bool {
	return len(k) == 0
}

type Music struct {
	date     time.Time
	title    string
	length   time.Duration
	platform string
	id       string
}

func (m *Music) Key() MusicKey {
	return MusicKey(m.platform + m.id)
}

func (m *Music) Date() time.Time {
	return m.date
}

func (m Music) Title() string {
	return m.title[:len(m.title)-4] //remove .mp3
}

func (m *Music) Length() time.Duration {
	return m.length
}

type musicJson struct {
	Date     int64  `json:"date"` //epoch seconds
	Title    string `json:"title"`
	Length   string `json:"length"`
	Platform string `json:"platform"`
	ID       string `json:"id"`
}

func (m *Music) MarshalJSON() ([]byte, error) {
	return json.Marshal(musicJson{m.date.Unix(), m.title, m.length.String(), m.platform, m.id})
}

func (m *Music) UnmarshalJSON(data []byte) error {
	var jsonData musicJson
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	length, err := time.ParseDuration(jsonData.Length)
	if err != nil {
		return err
	}

	*m = Music{time.Unix(jsonData.Date, 0), jsonData.Title, length, jsonData.Platform, jsonData.ID}
	return nil
}

type FileSystem interface {
	initialize() error

	//Get all the album keys from the file system.
	getAllAlbums() ([]Album, error)

	//Get the album from the file system by key.
	getAlbum(AlbumKey) (Album, error)

	//Create the album in the file system. Return the generated key.
	createAlbum(Album) (AlbumKey, error)

	//Update the album in the file system.
	//
	//The album must already exist in the system by its key.
	updateAlbum(Album) error

	//Remove the album from the file system by key.
	removeAlbum(AlbumKey) error

	//Get music from the file system by key.
	getMusic(MusicKey) (io.ReadSeekCloser, error)

	//Upload the music to the file system. Return the generated key.
	uploadMusic(Music, io.Reader) (MusicKey, error)

	//Remove the music from the file system by key.
	removeMusic(MusicKey) error
}
