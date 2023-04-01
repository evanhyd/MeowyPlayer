package cwidget

import (
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/player"
)

type albumNameFilter struct {
	itemList *AlbumItemList
}

func (filter *albumNameFilter) Notify(text string) {
	lowerCaseText := strings.ToLower(text)
	filter.itemList.SetFilter(func(album player.Album) bool {
		return strings.Contains(strings.ToLower(album.Title()), lowerCaseText)
	})
	filter.itemList.ScrollToTop()
}

type albumNameSorter struct {
	itemList *AlbumItemList
	reverse  bool
}

func (sorter *albumNameSorter) Notify() {
	sorter.itemList.SetSorter(func(album0, album1 player.Album) bool {
		return (strings.Compare(strings.ToLower(album0.Title()), strings.ToLower(album1.Title())) < 0) != sorter.reverse
	})
	sorter.reverse = !sorter.reverse
}

type albumDateSorter struct {
	itemList *AlbumItemList
	reverse  bool
}

func (sorter *albumDateSorter) Notify() {
	sorter.itemList.SetSorter(func(album0, album1 player.Album) bool {
		return (album0.ModifiedDate().Compare(album1.ModifiedDate()) > 0) != sorter.reverse
	})
	sorter.reverse = !sorter.reverse
}

type albumItemUpdater struct {
	itemList *AlbumItemList
}

func (itemUpdater *albumItemUpdater) Notify(albums []player.Album) {
	itemUpdater.itemList.SetItems(albums)
}

type AlbumItemList struct {
	ItemList[player.Album]
	nameFilter  albumNameFilter
	nameSorter  albumNameSorter
	dateSorter  albumDateSorter
	itemUpdater albumItemUpdater
}

func NewAlbumItemList(createItem func() fyne.CanvasObject, updateItem func(player.Album, fyne.CanvasObject)) *AlbumItemList {
	albumItemList := &AlbumItemList{}
	albumItemList.Initialize(createItem, updateItem)
	albumItemList.nameFilter = albumNameFilter{albumItemList}
	albumItemList.nameSorter = albumNameSorter{albumItemList, false}
	albumItemList.dateSorter = albumDateSorter{albumItemList, true}
	albumItemList.itemUpdater = albumItemUpdater{albumItemList}
	albumItemList.ExtendBaseWidget(albumItemList)
	return albumItemList
}

func (albumItemList *AlbumItemList) NameFilter() pattern.OneArgObserver[string] {
	return &albumItemList.nameFilter
}

func (albumItemList *AlbumItemList) NameSorter() pattern.ZeroArgObserver {
	return &albumItemList.nameSorter
}

func (albumItemList *AlbumItemList) DateFilter() pattern.ZeroArgObserver {
	return &albumItemList.dateSorter
}

func (albumItemList *AlbumItemList) ItemUpdater() pattern.OneArgObserver[[]player.Album] {
	return &albumItemList.itemUpdater
}
