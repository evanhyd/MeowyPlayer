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
	const musicTabName = "Music"

	//music views
	data := cbinding.MakeDataList[player.Music]()
	view := newMusicViewList(&data)
	manager.GetCurrentAlbum().Attach(utility.MakeCallback(func(a *player.Album) { data.Notify(a.MusicList) }))

	border := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, newMusicTitleButton(&data, view), newMusicDateButton(&data, view)),
			nil,
			container.NewGridWithRows(1, newMusicAdderLocalButton(&data, view), newMusicAdderOnlineButton(&data, view)),
			newMusicSearchBar(&data, view),
		),
		nil,
		nil,
		nil,
		view,
	)
	return container.NewTabItemWithIcon(musicTabName, resource.MusicTabIcon(), border)
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
		data.SetFilter(func(a player.Music) bool {
			return strings.Contains(strings.ToLower(a.Title), title)
		})
	}
	return entry
}

func newMusicAdderLocalButton(data *cbinding.DataList[player.Music], view *cwidget.MusicViewList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.MusicAdderLocalIcon(), func() {
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
	button := widget.NewButtonWithIcon("", resource.MusicAdderOnlineIcon(), func() {
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
		data.SetSorter(func(a1, a2 player.Music) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.Importance = widget.LowImportance
	button.OnTapped()
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
