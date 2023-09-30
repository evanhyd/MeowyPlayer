package cwidget

import (
	"strings"

	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/ui/cbinding"
)

// make data sort by music title
func NewMusicTitleButton(data *cbinding.MusicDataList, title string) *widget.Button {
	reverse := false
	button := widget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Music) bool {
			return (strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) < 0) != reverse
		})
	})
	button.Importance = widget.LowImportance
	return button
}

// make data sort by music date
func NewMusicDateButton(data *cbinding.MusicDataList, title string) *widget.Button {
	reverse := true
	button := widget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Music) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.Importance = widget.LowImportance
	button.OnTapped()
	return button
}

func NewAlbumTitleButton(data *cbinding.AlbumDataList, title string) *widget.Button {
	reverse := false
	button := widget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Album) bool {
			return (strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) < 0) != reverse
		})
	})
	button.Importance = widget.LowImportance
	return button
}

func NewAlbumDateButton(data *cbinding.AlbumDataList, title string) *widget.Button {
	reverse := true
	button := widget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Album) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.Importance = widget.LowImportance
	button.OnTapped()
	return button
}
