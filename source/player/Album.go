package player

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
)

type Album struct {
	Date      time.Time     `json:"date"`
	Title     string        `json:"title"`
	MusicList []Music       `json:"musicList"`
	Cover     fyne.Resource `json:"-"`
}

func (a *Album) Description() string {
	return fmt.Sprintf("[%v]\n%v", len(a.MusicList), a.Date.Format(time.DateTime))
}
