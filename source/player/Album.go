package player

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/canvas"
)

type Album struct {
	title        string
	modifiedDate time.Time
	musicNumber  int
	coverIcon    *canvas.Image
}

func (album *Album) Title() string {
	return album.title
}

func (album *Album) ModifiedDate() time.Time {
	return album.modifiedDate
}

func (album *Album) CoverIcon() *canvas.Image {
	return album.coverIcon
}

func (album *Album) Description() string {
	year, month, day := album.modifiedDate.Date()
	return fmt.Sprintf("%v\nMusic: %v\nModified: %v %v %v", album.title, album.musicNumber, year, month, day)
}

func (album *Album) IsPlaceHolder() bool {
	return album.title == ""
}

func GetPlaceHolderAlbum() Album {
	return Album{}
}
