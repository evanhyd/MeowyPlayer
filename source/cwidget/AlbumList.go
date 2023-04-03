package cwidget

import (
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/player"
)

type AlbumList struct {
	List[player.Album]

	reverseTitle bool
	reverseDate  bool
}

func NewAlbumList(createItem func() fyne.CanvasObject, updateItem func(player.Album, fyne.CanvasObject)) *AlbumList {
	albumList := &AlbumList{}
	albumList.Initialize(createItem, updateItem)
	albumList.ExtendBaseWidget(albumList)
	albumList.reverseDate = true
	return albumList
}

// func (albumItemList *AlbumItemList) Secondary {
// albumItemList.ItemList.List
// }

func (albumList *AlbumList) SetTitleFilter(title string) {
	lowerCaseTitle := strings.ToLower(title)
	albumList.SetFilter(func(album player.Album) bool {
		return strings.Contains(strings.ToLower(album.Title()), lowerCaseTitle)
	})
	albumList.ScrollToTop()
}

func (albumList *AlbumList) SetTitleSorter() {
	albumList.SetSorter(func(album0, album1 player.Album) bool {
		return (strings.Compare(strings.ToLower(album0.Title()), strings.ToLower(album1.Title())) < 0) != albumList.reverseTitle
	})
	albumList.reverseTitle = !albumList.reverseTitle
	albumList.reverseDate = false
}

func (albumList *AlbumList) SetDateSorter() {
	albumList.SetSorter(func(album0, album1 player.Album) bool {
		return (album0.ModifiedDate().Compare(album1.ModifiedDate()) > 0) != albumList.reverseDate
	})
	albumList.reverseDate = !albumList.reverseDate
	albumList.reverseTitle = false
}
