package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cbinding"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newAlbumTab() *container.TabItem {
	var selectedAlbum resource.Album
	data := cbinding.MakeDataList[resource.Album]()
	client.Manager().AddCollectionListener(pattern.MakeCallback(func(c resource.Collection) {
		data.Notify(maps.Values(c.Albums))
	}))

	// renaming title dialog
	renameEntry := widget.NewEntry()
	renameDialog := dialog.NewCustomConfirm("Enter title:", "Confirm", "Cancel", renameEntry, func(confirm bool) {
		if confirm {
			showErrorIfAny(client.Manager().RenameAlbum(selectedAlbum, renameEntry.Text))
		}
	}, getWindow())

	// updating album icon dialog
	coverDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			showErrorIfAny(err)
		} else if result != nil {
			showErrorIfAny(client.Manager().EditCover(selectedAlbum, result.URI().Path()))
		}
	}, getWindow())
	coverDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
	coverDialog.SetConfirmText("Upload")

	// deleting album dialog
	deleteDialog := dialog.NewConfirm("", "Do you want to delete the album?", func(delete bool) {
		if delete {
			showErrorIfAny(client.Manager().DeleteAlbum(selectedAlbum))
		}
	}, getWindow())

	// pop up menu
	editingMenu := widget.NewPopUpMenu(fyne.NewMenu("",
		fyne.NewMenuItem("Rename", func() {
			renameEntry.SetText("")
			renameDialog.Show()
			renameEntry.FocusGained()
		}),
		fyne.NewMenuItem("Edit Cover", coverDialog.Show),
		fyne.NewMenuItem("Delete", deleteDialog.Show)), getWindow().Canvas())

	// album views
	albumViews := cwidget.NewViewList(&data, container.NewGridWrap(fyne.NewSize(135.0, 165.0)),
		func(album resource.Album) fyne.CanvasObject {
			view := cwidget.NewAlbumView(&album)
			view.OnTapped = func(*fyne.PointEvent) {
				showErrorIfAny(client.Manager().SetFocusedAlbum(album))
			}
			view.OnTappedSecondary = func(event *fyne.PointEvent) {
				selectedAlbum = album
				editingMenu.ShowAtPosition(event.AbsolutePosition)
			}
			return view
		},
	)

	// title search bar
	searchBar := widget.NewEntry()
	searchBar.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(album resource.Album) bool {
			return strings.Contains(strings.ToLower(album.Title), title)
		})
	}

	// title sorting button
	ascendTitle := -1
	titleButton := cwidget.NewButton("Title", func() {
		ascendTitle = -ascendTitle
		data.SetSorter(func(a1, a2 resource.Album) int {
			return strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) * ascendTitle
		})
	})

	// date sorting button
	ascendDate := 1
	dateButton := cwidget.NewButton("Date", func() {
		ascendDate = -ascendDate
		data.SetSorter(func(a1, a2 resource.Album) int {
			return a1.Date.Compare(a2.Date) * ascendDate
		})
	})
	defer dateButton.OnTapped()

	// add album button
	addAlbumButton := cwidget.NewButtonWithIcon("", theme.ContentAddIcon(), func() { showErrorIfAny(client.AddRandomAlbum()) })

	//progress bar
	progressBar := widget.NewProgressBar()
	progressDialog := dialog.NewCustomWithoutButtons("loading", progressBar, getWindow())
	syncMusicButton := cwidget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		progressDialog.Show()
		remains := client.SyncCollection()
		for remain := range remains {
			progressBar.SetValue(remain)
		}
		progressDialog.Hide()
	})

	return container.NewTabItemWithIcon("Album", resource.AlbumTabIcon, container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, titleButton, dateButton),
			nil,
			container.NewGridWithRows(1, addAlbumButton, syncMusicButton),
			searchBar,
		),
		nil,
		nil,
		nil,
		albumViews,
	))
}
