package cwidget

import (
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/player"
)

type MusicItemList struct {
	ItemList[player.Music]

	reverseTitle bool
	reverseDate  bool
}

func NewMusicItemList(createItem func() fyne.CanvasObject, updateItem func(player.Music, fyne.CanvasObject)) *MusicItemList {
	musicItemList := &MusicItemList{}
	musicItemList.Initialize(createItem, updateItem)
	musicItemList.ExtendBaseWidget(musicItemList)
	musicItemList.reverseDate = true
	return musicItemList
}

func (musicItemList *MusicItemList) SetTitleFilter(title string) {
	lowerCaseTitle := strings.ToLower(title)
	musicItemList.SetFilter(func(music player.Music) bool {
		return strings.Contains(strings.ToLower(music.Title()), lowerCaseTitle)
	})
	musicItemList.ScrollToTop()
}

func (musicItemList *MusicItemList) SetTitleSorter() {
	musicItemList.SetSorter(func(music0, music1 player.Music) bool {
		return (strings.Compare(strings.ToLower(music0.Title()), strings.ToLower(music1.Title())) < 0) != musicItemList.reverseTitle
	})
	musicItemList.reverseTitle = !musicItemList.reverseTitle
	musicItemList.reverseDate = false
}

func (musicItemList *MusicItemList) SetDateSorter() {
	musicItemList.SetSorter(func(music0, music1 player.Music) bool {
		return (music0.ModifiedDate().Compare(music1.ModifiedDate()) > 0) != musicItemList.reverseDate
	})
	musicItemList.reverseDate = !musicItemList.reverseDate
	musicItemList.reverseTitle = false
}
