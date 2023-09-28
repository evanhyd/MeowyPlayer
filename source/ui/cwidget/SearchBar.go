package cwidget

import (
	"strings"

	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui/cbinding"
)

func NewAlbumSearchBar(data *cbinding.AlbumDataList) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(a player.Album) bool {
			return strings.Contains(strings.ToLower(a.Title), title)
		})
	}
	return entry
}

func NewMusicSearchBar(data *cbinding.MusicDataList) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(a player.Music) bool {
			return strings.Contains(strings.ToLower(a.SimpleTitle()), title)
		})
	}
	return entry
}
