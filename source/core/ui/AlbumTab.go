package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cbinding"
	"meowyplayer.com/core/ui/cwidget"
)

func newAlbumTab() *container.TabItem {
	makeRenameDialog := func(album *resource.Album) func() {
		entry := widget.NewEntry()
		return func() {
			dialog.ShowCustomConfirm("Enter title:", "Confirm", "Cancel", entry, func(rename bool) {
				if rename {
					showErrorIfAny(client.Manager().UpdateAlbumTitle(*album, entry.Text))
				}
			}, getWindow())
		}
	}

	makeCoverDialog := func(album *resource.Album) func() {
		return func() {
			fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
				if err != nil {
					showErrorIfAny(err)
				} else if result != nil {
					showErrorIfAny(client.Manager().UpdateAlbumCover(*album, result.URI().Path()))
				}
			}, getWindow())
			fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
			fileOpenDialog.SetConfirmText("Upload")
			fileOpenDialog.Show()
		}
	}

	makeDeleteAlbumDialog := func(album *resource.Album) func() {
		return func() {
			dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", album.Title), func(delete bool) {
				if delete {
					showErrorIfAny(client.Manager().DeleteAlbum(*album))
				}
			}, getWindow())
		}
	}

	showAlbumMenu := func(album *resource.Album, canvas fyne.Canvas, pos fyne.Position) {
		rename := fyne.NewMenuItem("Rename", makeRenameDialog(album))
		cover := fyne.NewMenuItem("Cover", makeCoverDialog(album))
		delete := fyne.NewMenuItem("Delete", makeDeleteAlbumDialog(album))
		widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", rename, cover, delete), canvas, pos)
	}

	newAlbumViewList := func(data *cbinding.AlbumDataList) *cwidget.ViewList[resource.Album] {
		return cwidget.NewViewList(data, container.NewGridWrap(fyne.NewSize(135.0, 165.0)),
			func(album resource.Album) fyne.CanvasObject {
				view := cwidget.NewAlbumView(&album)
				view.OnTapped = func(*fyne.PointEvent) {
					client.Manager().SetAlbum(album)
				}
				view.OnTappedSecondary = func(event *fyne.PointEvent) {
					canvas := fyne.CurrentApp().Driver().CanvasForObject(view)
					showAlbumMenu(&album, canvas, event.AbsolutePosition)
				}
				return view
			},
		)
	}

	newAlbumSearchBar := func(data *cbinding.AlbumDataList) *widget.Entry {
		entry := widget.NewEntry()
		entry.OnChanged = func(title string) {
			title = strings.ToLower(title)
			data.SetFilter(func(a resource.Album) bool {
				return strings.Contains(strings.ToLower(a.Title), title)
			})
		}
		return entry
	}

	newAlbumTitleButton := func(data *cbinding.AlbumDataList, title string) *widget.Button {
		order := -1
		return cwidget.NewButton(title, func() {
			order = -order
			data.SetSorter(func(a1, a2 resource.Album) int {
				return strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) * order
			})
		})
	}

	newAlbumDateButton := func(data *cbinding.AlbumDataList, title string) *widget.Button {
		order := 1
		button := cwidget.NewButton(title, func() {
			order = -order
			data.SetSorter(func(a1, a2 resource.Album) int {
				return a1.Date.Compare(a2.Date) * order
			})
		})
		button.OnTapped()
		return button
	}

	data := cbinding.MakeAlbumDataList()
	client.Manager().AddCollectionListener(&data)

	albumAdderLocal := cwidget.NewButtonWithIcon("", theme.ContentAddIcon(), showAddLocalAlbumDialog)
	albumAdderOnline := cwidget.NewButtonWithIcon("", resource.AlbumAdderOnlineIcon, showAddOnlineAlbumDialog)

	return container.NewTabItemWithIcon("Album", resource.AlbumTabIcon, container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, newAlbumTitleButton(&data, "Title"), newAlbumDateButton(&data, "Date")),
			nil,
			container.NewGridWithRows(1, albumAdderLocal, albumAdderOnline),
			newAlbumSearchBar(&data),
		),
		nil,
		nil,
		nil,
		newAlbumViewList(&data),
	))
}
