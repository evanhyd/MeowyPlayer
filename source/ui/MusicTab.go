package ui

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newMusicTab() *container.TabItem {
	data := cbinding.MakeMusicDataList()
	client.GetAlbumData().Attach(&data)

	searchBar := newMusicSearchBar(&data)
	client.GetAlbumData().Attach(pattern.MakeCallback(func(*resource.Album) { searchBar.SetText("") }))

	musicAdderLocal := cwidget.NewButtonWithIcon("", theme.FolderOpenIcon(), showAddLocalMusicDialog)
	musicAdderOnline := cwidget.NewButtonWithIcon("", resource.MusicAdderOnlineIcon(), showAddOnlineMusicDialog)

	return container.NewTabItemWithIcon("Music", resource.MusicTabIcon(), container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, newMusicTitleButton(&data, "Title"), newMusicDateButton(&data, "Date")),
			nil,
			container.NewGridWithRows(1, musicAdderLocal, musicAdderOnline),
			searchBar,
		),
		nil,
		nil,
		nil,
		newMusicViewList(&data),
	))
}

func newMusicViewList(data *cbinding.MusicDataList) *cwidget.ViewList[resource.Music] {
	return cwidget.NewViewList[resource.Music](data, container.NewVBox(),
		func(music resource.Music) fyne.CanvasObject {
			view := cwidget.NewMusicView(&music)
			view.OnTapped = func(*fyne.PointEvent) { client.GetPlayListData().Set(resource.NewPlayList(data.GetAlbum(), &music)) }
			view.OnTappedSecondary = func(*fyne.PointEvent) { showDeleteMusicDialog(&music) }
			return view
		},
	)
}

func showDeleteMusicDialog(music *resource.Music) {
	dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", music.Title), func(delete bool) {
		if delete {
			log.Printf("delete %v from the album %v\n", music.Title, client.GetAlbumData().Get().Title)
			showErrorIfAny(client.DeleteMusic(music))
		}
	}, getWindow())
}

func newMusicSearchBar(data *cbinding.MusicDataList) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(a resource.Music) bool {
			return strings.Contains(strings.ToLower(a.SimpleTitle()), title)
		})
	}
	return entry
}

// make data sort by music title
func newMusicTitleButton(data *cbinding.MusicDataList, title string) *widget.Button {
	reverse := false
	return cwidget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 resource.Music) bool {
			return (strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) < 0) != reverse
		})
	})
}

// make data sort by music date
func newMusicDateButton(data *cbinding.MusicDataList, title string) *widget.Button {
	reverse := true
	button := cwidget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 resource.Music) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.OnTapped()
	return button
}
