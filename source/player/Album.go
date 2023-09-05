package player

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/utility"
)

type Album struct {
	utility.Subject[*Album]
	Date      time.Time     `json:"date"`
	Title     string        `json:"title"`
	MusicList []Music       `json:"musicList"`
	Cover     *canvas.Image `json:"-"`
}

func (a *Album) Description() string {
	return fmt.Sprintf("[%v] %v\n%v", len(a.MusicList), a.Title, a.Date.Format(time.DateTime))
}

func (a *Album) NotifyAll() {
	a.Subject.NotifyAll(a)
}
