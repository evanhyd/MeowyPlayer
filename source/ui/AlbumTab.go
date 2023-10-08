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
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/ui/cwidget"
)

func newAlbumTab() *container.TabItem {
	data := cbinding.MakeAlbumDataList()
	client.GetCollectionData().Attach(&data)

	albumAdderLocal := cwidget.NewButtonWithIcon("", resource.AlbumAdderLocalIcon(), showAddLocalAlbumDialog)
	albumAdderOnline := cwidget.NewButtonWithIcon("", resource.AlbumAdderOnlineIcon(), showAddOnlineAlbumDialog)

	return container.NewTabItemWithIcon("Album", resource.AlbumTabIcon(), container.NewBorder(
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

func newAlbumViewList(data *cbinding.AlbumDataList) *cwidget.AlbumViewList {
	return cwidget.NewAlbumViewList(data, func(album *resource.Album) fyne.CanvasObject {
		view := cwidget.NewAlbumView(album)
		view.OnTapped = func(*fyne.PointEvent) { client.GetAlbumData().Set(album) }
		view.OnTappedSecondary = func(event *fyne.PointEvent) {
			canvas := fyne.CurrentApp().Driver().CanvasForObject(view)
			newAlbumMenu(canvas, album).ShowAtPosition(event.AbsolutePosition)
		}
		return view
	}, fyne.NewSize(135.0, 165.0))
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

func newAlbumMenu(canvas fyne.Canvas, album *resource.Album) *widget.PopUpMenu {
	rename := fyne.NewMenuItem("Rename", makeRenameDialog(album))
	cover := fyne.NewMenuItem("Cover", makeCoverDialog(album))
	delete := fyne.NewMenuItem("Delete", makeDeleteAlbumDialog(album))
	return widget.NewPopUpMenu(fyne.NewMenu("", rename, cover, delete), canvas)
}

func makeRenameDialog(album *resource.Album) func() {
	entry := widget.NewEntry()
	return func() {
		dialog.ShowCustomConfirm("Enter title:", "Confirm", "Cancel", entry, func(rename bool) {
			if rename {
				log.Printf("rename %v to %v\n", album.Title, entry.Text)
				showErrorIfAny(client.UpdateAlbumTitle(album, entry.Text))
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
				log.Printf("update %v's cover: %v\n", album.Title, result.URI().Path())
				showErrorIfAny(client.UpdateAlbumCover(album, result.URI().Path()))
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
				log.Printf("delete %v\n", album.Title)
				showErrorIfAny(client.DeleteAlbum(album))
			}
		}, getWindow())
	}
}
