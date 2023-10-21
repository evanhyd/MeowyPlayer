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
	data := cbinding.MakeAlbumDataList()
	client.GetInstance().AddCollectionListener(&data)

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

func newAlbumViewList(data *cbinding.AlbumDataList) *cwidget.ViewList[resource.Album] {
	return cwidget.NewViewList(data, container.NewGridWrap(fyne.NewSize(135.0, 165.0)),
		func(album resource.Album) fyne.CanvasObject {
			view := cwidget.NewAlbumView(&album)
			view.OnTapped = func(*fyne.PointEvent) { client.GetInstance().SetAlbum(album) }
			view.OnTappedSecondary = func(event *fyne.PointEvent) {
				canvas := fyne.CurrentApp().Driver().CanvasForObject(view)
				showAlbumMenu(&album, canvas, event.AbsolutePosition)
			}
			return view
		},
	)
}

func newAlbumSearchBar(data *cbinding.AlbumDataList) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(a resource.Album) bool {
			return strings.Contains(strings.ToLower(a.Title), title)
		})
	}
	return entry
}

func newAlbumTitleButton(data *cbinding.AlbumDataList, title string) *widget.Button {
	reverse := false
	return cwidget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 resource.Album) bool {
			return (strings.ToLower(a1.Title) < strings.ToLower(a2.Title)) != reverse
		})
	})
}

func newAlbumDateButton(data *cbinding.AlbumDataList, title string) *widget.Button {
	reverse := true
	button := cwidget.NewButton(title, func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 resource.Album) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.OnTapped()
	return button
}

func showAlbumMenu(album *resource.Album, canvas fyne.Canvas, pos fyne.Position) {
	rename := fyne.NewMenuItem("Rename", makeRenameDialog(album))
	cover := fyne.NewMenuItem("Cover", makeCoverDialog(album))
	delete := fyne.NewMenuItem("Delete", makeDeleteAlbumDialog(album))
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", rename, cover, delete), canvas, pos)
}

func makeRenameDialog(album *resource.Album) func() {
	entry := widget.NewEntry()
	return func() {
		dialog.ShowCustomConfirm("Enter title:", "Confirm", "Cancel", entry, func(rename bool) {
			if rename {
				showErrorIfAny(client.GetInstance().UpdateAlbumTitle(*album, entry.Text))
			}
		}, getWindow())
	}
}

func makeCoverDialog(album *resource.Album) func() {
	return func() {
		fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorIfAny(err)
			} else if result != nil {
				showErrorIfAny(client.GetInstance().UpdateAlbumCover(*album, result.URI().Path()))
			}
		}, getWindow())
		fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
		fileOpenDialog.SetConfirmText("Upload")
		fileOpenDialog.Show()
	}
}

func makeDeleteAlbumDialog(album *resource.Album) func() {
	return func() {
		dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", album.Title), func(delete bool) {
			if delete {
				showErrorIfAny(client.GetInstance().DeleteAlbum(*album))
			}
		}, getWindow())
	}
}
