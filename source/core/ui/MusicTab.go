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
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/player"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cbinding"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newMusicTab() *container.TabItem {
	newMusicSearchBar := func(data *cbinding.MusicDataList) *widget.Entry {
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
	newMusicTitleButton := func(data *cbinding.MusicDataList, title string) *widget.Button {
		reverse := -1
		return cwidget.NewButton(title, func() {
			reverse = -reverse
			data.SetSorter(func(a1, a2 resource.Music) int {
				return strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) * reverse
			})
		})
	}

	// make data sort by music date
	newMusicDateButton := func(data *cbinding.MusicDataList, title string) *widget.Button {
		reverse := 1
		button := cwidget.NewButton(title, func() {
			reverse = -reverse
			data.SetSorter(func(a1, a2 resource.Music) int {
				return a1.Date.Compare(a2.Date) * reverse
			})
		})
		button.OnTapped()
		return button
	}

	showDeleteMusicDialog := func(music *resource.Music) {
		dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", music.Title), func(delete bool) {
			if delete {
				log.Printf("delete %v from the album %v\n", music.Title, client.Manager().Album().Title)
				showErrorIfAny(client.Manager().DeleteMusic(*music))
			}
		}, getWindow())
	}

	newMusicViewList := func(data *cbinding.MusicDataList) *cwidget.ViewList[resource.Music] {
		return cwidget.NewViewList(data, container.NewVBox(),
			func(music resource.Music) fyne.CanvasObject {
				view := cwidget.NewMusicView(&music)
				view.OnTapped = func(*fyne.PointEvent) { client.Manager().SetPlayList(player.NewPlayList(data.MusicList(), &music)) }
				view.OnTappedSecondary = func(*fyne.PointEvent) { showDeleteMusicDialog(&music) }
				return view
			},
		)
	}

	data := cbinding.MakeMusicDataList()
	searchBar := newMusicSearchBar(&data)
	client.Manager().AddAlbumListener(&data)
	client.Manager().AddAlbumListener(pattern.MakeCallback(func(resource.Album) { searchBar.SetText("") }))
	musicAdderLocal := cwidget.NewButtonWithIcon("", theme.FolderOpenIcon(), showAddLocalMusicDialog)
	musicAdderOnline := cwidget.NewButtonWithIcon("", resource.MusicAdderOnlineIcon, showAddOnlineMusicDialog)

	return container.NewTabItemWithIcon("Music", resource.MusicTabIcon, container.NewBorder(
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