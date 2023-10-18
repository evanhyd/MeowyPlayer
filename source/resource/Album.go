package resource

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
)

type Album struct {
	Date      time.Time        `json:"date"`
	Title     string           `json:"title"`
	MusicList map[string]Music `json:"musicList"`
	Cover     fyne.Resource    `json:"-"`
}

func (a *Album) Description() string {
	return fmt.Sprintf("%v\n\nMusic: %v\n\n%v", a.Title, len(a.MusicList), a.Date.Format(time.DateTime))
}
