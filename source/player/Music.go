package player

import (
	"fmt"
	"time"
)

type Music struct {
	title        string
	duration     time.Duration
	modifiedDate time.Time
}

func (music *Music) Title() string {
	return music.title
}

func (music *Music) Duration() time.Duration {
	return music.duration
}

func (music *Music) ModifiedDate() time.Time {
	return music.modifiedDate
}

func (music *Music) Description() string {
	return fmt.Sprintf("%02d:%02d", int(music.duration.Minutes())%60, int(music.duration.Seconds())%60) + " | " + music.title
}
