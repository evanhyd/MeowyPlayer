package cwidget

import (
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/player"
)

type AlbumItemList struct {
	ItemList[player.Album]

	reverseTitle bool
	reverseDate  bool
}

func NewAlbumItemList(createItem func() fyne.CanvasObject, updateItem func(player.Album, fyne.CanvasObject)) *AlbumItemList {
	albumItemList := &AlbumItemList{}
	albumItemList.Initialize(createItem, updateItem)
	albumItemList.ExtendBaseWidget(albumItemList)
	albumItemList.reverseDate = true
	return albumItemList
}

func (albumItemList *AlbumItemList) SetTitleFilter(title string) {
	lowerCaseTitle := strings.ToLower(title)
	albumItemList.SetFilter(func(album player.Album) bool {
		return strings.Contains(strings.ToLower(album.Title()), lowerCaseTitle)
	})
	albumItemList.ScrollToTop()
}

func (albumItemList *AlbumItemList) SetTitleSorter() {
	albumItemList.SetSorter(func(album0, album1 player.Album) bool {
		return (strings.Compare(strings.ToLower(album0.Title()), strings.ToLower(album1.Title())) < 0) != albumItemList.reverseTitle
	})
	albumItemList.reverseTitle = !albumItemList.reverseTitle
	albumItemList.reverseDate = false
}

func (albumItemList *AlbumItemList) SetDateSorter() {
	albumItemList.SetSorter(func(album0, album1 player.Album) bool {
		return (album0.ModifiedDate().Compare(album1.ModifiedDate()) > 0) != albumItemList.reverseDate
	})
	albumItemList.reverseDate = !albumItemList.reverseDate
	albumItemList.reverseTitle = false
}
