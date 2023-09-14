package ui

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/source/utility"
)

func newMusicTab() *container.TabItem {
	const (
		musicTabName     = "Music"
		musicTabIconName = "music_tab.png"
	)

	//music views
	data := cbinding.MakeDataList[player.Music]()
	view := newMusicViewList(&data)
	manager.GetCurrentAlbum().Attach(utility.MakeCallback(func(a *player.Album) { data.Notify(a.MusicList) }))

	searchBar := newMusicSearchBar(&data, view)
	musicAdderLocalButton := newMusicAdderLocalButton(&data, view)
	musicAdderOnlineButton := newMusicAdderOnlineButton(&data, view)
	titleButton := newMusicTitleButton(&data, view)
	dateButton := newMusicDateButton(&data, view)
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

	return container.NewTabItemWithIcon(musicTabName, resource.GetAsset(musicTabIconName), border)
}

func newMusicViewList(data *cbinding.DataList[player.Music]) *cwidget.MusicViewList {
	list := cwidget.NewMusicViewList(func(m *player.Music) fyne.CanvasObject {
		view := cwidget.NewMusicView(m)
		view.OnTapped = func(*fyne.PointEvent) { fmt.Println(m.Title) }
		view.OnTappedSecondary = func(*fyne.PointEvent) { showDeleteMusicDialog(m) }
		return view
	})

	data.Attach(list)

	return list
}

func newMusicSearchBar(data *cbinding.DataList[player.Music], view *cwidget.MusicViewList) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		filter := func(a player.Music) bool { return strings.Contains(strings.ToLower(a.Title), title) }
		data.SetFilter(filter)
	}
	return entry
}

func newMusicAdderLocalButton(data *cbinding.DataList[player.Music], view *cwidget.MusicViewList) *widget.Button {
	const iconName = "music_adder_local.png"
	button := widget.NewButtonWithIcon("", resource.GetAsset(iconName), func() {
		fileReader := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorIfAny(err)
			} else if result != nil {
				log.Printf("add %v from local to %v\n", result.URI().Name(), manager.GetCurrentAlbum().Get().Title)
				showErrorIfAny(manager.AddMusic(result))
			}
		}, getMainWindow())
		fileReader.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
		fileReader.SetConfirmText("Add")
		fileReader.Show()
	})
	button.Importance = widget.LowImportance

	return button
}

func newMusicAdderOnlineButton(data *cbinding.DataList[player.Music], view *cwidget.MusicViewList) *widget.Button {
	const iconName = "music_adder_online.png"
	const musicAdderOnlineSearchIconName = "music_adder_online_search.png" //move to other place
	button := widget.NewButtonWithIcon("", resource.GetAsset(iconName), func() {
		log.Println("add music from online")
		//to do
	})
	button.Importance = widget.LowImportance
	return button
}

func newMusicTitleButton(data *cbinding.DataList[player.Music], view *cwidget.MusicViewList) *widget.Button {
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

func newMusicDateButton(data *cbinding.DataList[player.Music], view *cwidget.MusicViewList) *widget.Button {
	reverse := true
	button := widget.NewButton("Date", func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Music) bool { return a1.Date.After(a2.Date) != reverse })
	})
	button.Importance = widget.LowImportance
	return button
}

func showDeleteMusicDialog(music *player.Music) {
	dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", music.Title), func(delete bool) {
		if delete {
			log.Printf("delete %vfrom the album %v \n", music.Title, manager.GetCurrentAlbum().Get().Title)
			showErrorIfAny(manager.DeleteMusic(music))
		}
	}, getMainWindow())
}
