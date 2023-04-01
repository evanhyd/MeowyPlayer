package cwidget

import (
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/player"
)

type musicNameFilter struct {
	itemList *MusicItemList
}

func (filter *musicNameFilter) Notify(text string) {
	lowerCaseText := strings.ToLower(text)
	filter.itemList.SetFilter(func(music player.Music) bool {
		return strings.Contains(strings.ToLower(music.Title()), lowerCaseText)
	})
	filter.itemList.ScrollToTop()
}

type musicNameSorter struct {
	itemList *MusicItemList
	reverse  bool
}

func (sorter *musicNameSorter) Notify() {
	sorter.itemList.SetSorter(func(music0, music1 player.Music) bool {
		return (strings.Compare(strings.ToLower(music0.Title()), strings.ToLower(music1.Title())) < 0) != sorter.reverse
	})
	sorter.reverse = !sorter.reverse
}

type musicDateSorter struct {
	itemList *MusicItemList
	reverse  bool
}

func (sorter *musicDateSorter) Notify() {
	sorter.itemList.SetSorter(func(music0, music1 player.Music) bool {
		return (music0.ModifiedDate().Compare(music1.ModifiedDate()) > 0) != sorter.reverse
	})
	sorter.reverse = !sorter.reverse
}

type musicItemUpdater struct {
	itemList *MusicItemList
}

func (itemUpdater *musicItemUpdater) Notify(album player.Album, music []player.Music) {
	itemUpdater.itemList.SetItems(music)
	itemUpdater.itemList.ScrollToTop()
}

type MusicItemList struct {
	ItemList[player.Music]
	nameFilter  musicNameFilter
	nameSorter  musicNameSorter
	dateSorter  musicDateSorter
	itemUpdater musicItemUpdater
}

func NewMusicItemList(createItem func() fyne.CanvasObject, updateItem func(player.Music, fyne.CanvasObject)) *MusicItemList {
	musicItemList := &MusicItemList{}
	musicItemList.Initialize(createItem, updateItem)
	musicItemList.nameFilter = musicNameFilter{musicItemList}
	musicItemList.nameSorter = musicNameSorter{musicItemList, false}
	musicItemList.dateSorter = musicDateSorter{musicItemList, true}
	musicItemList.itemUpdater = musicItemUpdater{musicItemList}
	musicItemList.ExtendBaseWidget(musicItemList)
	return musicItemList
}

func (musicItemList *MusicItemList) NameFilter() pattern.OneArgObserver[string] {
	return &musicItemList.nameFilter
}

func (musicItemList *MusicItemList) NameSorter() pattern.ZeroArgObserver {
	return &musicItemList.nameSorter
}

func (musicItemList *MusicItemList) DateFilter() pattern.ZeroArgObserver {
	return &musicItemList.dateSorter
}

func (musicItemList *MusicItemList) ItemUpdater() pattern.TwoArgObserver[player.Album, []player.Music] {
	return &musicItemList.itemUpdater
}
