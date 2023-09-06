package player

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/canvas"
)

type Album struct {
	Date      time.Time     `json:"date"`
	Title     string        `json:"title"`
	MusicList []Music       `json:"musicList"`
	Cover     *canvas.Image `json:"-"`
}

func (a *Album) Description() string {
	return fmt.Sprintf("[%v] %v\n%v", len(a.MusicList), a.Title, a.Date.Format(time.DateTime))
}
