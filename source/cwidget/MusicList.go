package cwidget

import (
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/player"
)

type MusicList struct {
	List[player.Music]

	reverseTitle bool
	reverseDate  bool
}

func NewMusicList(createItem func() fyne.CanvasObject, updateItem func(player.Music, fyne.CanvasObject)) *MusicList {
	musicList := &MusicList{}
	musicList.Initialize(createItem, updateItem)
	musicList.ExtendBaseWidget(musicList)
	musicList.reverseDate = true
	return musicList
}

func (musicList *MusicList) SetTitleFilter(title string) {
	lowerCaseTitle := strings.ToLower(title)
	musicList.SetFilter(func(music *player.Music) bool {
		return strings.Contains(strings.ToLower(music.Title()), lowerCaseTitle)
	})
	musicList.ScrollToTop()
}

func (musicList *MusicList) SetTitleSorter() {
	musicList.SetSorter(func(music0, music1 *player.Music) bool {
		return (strings.Compare(strings.ToLower(music0.Title()), strings.ToLower(music1.Title())) < 0) != musicList.reverseTitle
	})
	musicList.reverseTitle = !musicList.reverseTitle
	musicList.reverseDate = false
}

func (musicList *MusicList) SetDateSorter() {
	musicList.SetSorter(func(music0, music1 *player.Music) bool {
		return (music0.ModifiedDate().Compare(music1.ModifiedDate()) > 0) != musicList.reverseDate
	})
	musicList.reverseDate = !musicList.reverseDate
	musicList.reverseTitle = false
}
