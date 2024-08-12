package model

import (
	"encoding/json"
	"time"
)

type MusicKey string

func (k MusicKey) IsEmpty() bool  { return len(k) == 0 }
func (k MusicKey) String() string { return string(k) }

type Music struct {
	date     time.Time
	title    string
	length   time.Duration
	platform string
	id       string
}

var _ json.Marshaler = Music{}
var _ json.Unmarshaler = &Music{}

func (m *Music) Key() MusicKey         { return MusicKey(m.platform + m.id) }
func (m *Music) Date() time.Time       { return m.date }
func (m Music) Title() string          { return m.title }
func (m *Music) Length() time.Duration { return m.length }

type musicJson struct {
	Date     int64  `json:"date"` //epoch seconds
	Title    string `json:"title"`
	Length   string `json:"length"`
	Platform string `json:"platform"`
	ID       string `json:"id"`
}

func (m Music) MarshalJSON() ([]byte, error) {
	return json.Marshal(musicJson{m.date.Unix(), m.title, m.length.String(), m.platform, m.id})
}

func (m *Music) UnmarshalJSON(data []byte) error {
	var buf musicJson
	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}
	length, err := time.ParseDuration(buf.Length)
	if err != nil {
		return err
	}

	*m = Music{time.Unix(buf.Date, 0), buf.Title, length, buf.Platform, buf.ID}
	return nil
}
