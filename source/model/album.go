package model

import (
	"encoding/json"
	"time"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

type AlbumKey struct{ id uuid.UUID }

var _ json.Marshaler = AlbumKey{}
var _ json.Unmarshaler = &AlbumKey{}

func newRandomAlbumKey() AlbumKey { return AlbumKey{uuid.New()} }
func (k AlbumKey) IsEmpty() bool  { return k.id == uuid.UUID{} }
func (k AlbumKey) String() string { return k.id.String() }

func (k AlbumKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *AlbumKey) UnmarshalJSON(d []byte) error {
	var err error
	k.id, err = uuid.ParseBytes(d)
	return err
}

type Album struct {
	key   AlbumKey
	date  time.Time
	title string
	music []Music
	cover fyne.Resource
}

var _ json.Marshaler = Album{}
var _ json.Unmarshaler = &Album{}

func (a *Album) Key() AlbumKey        { return a.key }
func (a *Album) Date() time.Time      { return a.date }
func (a Album) Title() string         { return a.title }
func (a *Album) Music() []Music       { return a.music }
func (a *Album) Cover() fyne.Resource { return a.cover }
func (a *Album) Count() int           { return len(a.music) }

type albumJson struct {
	Key   AlbumKey `json:"key"`
	Date  int64    `json:"date"` //epoch seconds
	Title string   `json:"title"`
	Music []Music  `json:"music"`
	Cover []byte   `json:"cover"`
}

func (a Album) MarshalJSON() ([]byte, error) {
	return json.Marshal(albumJson{a.key, a.date.Unix(), a.title, a.music, a.cover.Content()})
}

func (a *Album) UnmarshalJSON(data []byte) error {
	var buf albumJson
	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}
	*a = Album{buf.Key, time.Unix(buf.Date, 0), buf.Title, buf.Music, fyne.NewStaticResource(buf.Key.String(), buf.Cover)}
	return nil
}
