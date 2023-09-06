package ui

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource/album"
	"meowyplayer.com/source/resource/texture"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/utility"
)

func newMusicTab() *container.TabItem {
	const (
		musicTabName     = "Music"
		musicTabIconName = "music_tab.png"

		musicAdderOnlineSearchIconName = "music_adder_online_search.png" //move to other place
	)

	//music views
	data := cbinding.NewMusicList()
	view := newMusicView(data)
	album.Get().Attach(data)

	searchBar := newMusicSearchBar(data, view)
	musicAdderLocalButton := newMusicAdderLocalButton(data, view)
	musicAdderOnlineButton := newMusicAdderOnlineButton(data, view)
	titleButton := newMusicTitleButton(data, view)
	dateButton := newMusicDateButton(data, view)
	dateButton.OnTapped()

	border := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, titleButton, dateButton),
			nil,
			container.NewGridWithRows(1, musicAdderLocalButton, musicAdderOnlineButton),
			searchBar,
		),
		nil,
		nil,
		nil,
		view,
	)

	return container.NewTabItemWithIcon(musicTabName, texture.Get(musicTabIconName), border)
}

func newMusicView(data binding.DataList) *widget.List {
	view := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			setting := widget.NewButton("x", func() {})
			setting.Importance = widget.LowImportance
			intro := widget.NewLabel("")
			return container.NewBorder(nil, nil, setting, nil, intro)
		},
		func(item binding.DataItem, canvasObject fyne.CanvasObject) {
			data, err := item.(binding.Untyped).Get()
			utility.MustOk(err)
			music := data.(player.Music)

			objects := canvasObject.(*fyne.Container).Objects
			intro := objects[0].(*widget.Label)

			//optionally update
			if description := music.Description(); intro.Text != description {
				intro.Text = description

				//update setting functionality
				setting := objects[1].(*widget.Button)
				utility.MustNotNil(setting)
				setting.OnTapped = func() {
				}

				canvasObject.Refresh()
			}
		})

	//select and load album
	view.OnSelected = func(id widget.ListItemID) {
		item, err := data.GetItem(id)
		utility.MustOk(err)
		data, err := item.(binding.Untyped).Get()
		utility.MustOk(err)
		a := data.(player.Music)
		fmt.Println(a.Title)
		view.Unselect(id)
	}

	return view
}
func newMusicSearchBar(data *cbinding.MusicList, view *widget.List) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		filter := func(a player.Music) bool {
			return strings.Contains(strings.ToLower(a.Title), title)
		}
		data.SetFilter(filter)
		view.ScrollToTop()
	}
	return entry
}

func newMusicAdderLocalButton(data *cbinding.MusicList, view *widget.List) *widget.Button {
	const musicAdderLocalIconName = "music_adder_local.png"
	button := widget.NewButtonWithIcon("", texture.Get(musicAdderLocalIconName), func() {
		log.Println("add music from local")
		//to do
	})
	button.Importance = widget.LowImportance
	return button
}

func newMusicAdderOnlineButton(data *cbinding.MusicList, view *widget.List) *widget.Button {
	const musicAdderOnlineIconName = "music_adder_online.png"
	button := widget.NewButtonWithIcon("", texture.Get(musicAdderOnlineIconName), func() {
		log.Println("add music from online")
		//to do
	})
	button.Importance = widget.LowImportance
	return button
}

func newMusicTitleButton(data *cbinding.MusicList, view *widget.List) *widget.Button {
	reverse := false
	button := widget.NewButton("Title", func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Music) bool {
			return (strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) < 0) != reverse
		})
	})
	button.Importance = widget.LowImportance
	return button
}

func newMusicDateButton(data *cbinding.MusicList, view *widget.List) *widget.Button {
	reverse := true
	button := widget.NewButton("Date", func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Music) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.Importance = widget.LowImportance
	return button
}
