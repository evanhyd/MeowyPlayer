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

func newAlbumTab() *container.TabItem {
	const albumTabTitle = "Album"

	//bind data and view
	data := cbinding.MakeDataList[player.Album]()
	view := newAlbumViewList(&data)
	manager.GetCurrentConfig().Attach(utility.MakeCallback(func(config *player.Config) { data.Notify(config.Albums) }))

	border := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, newAlbumTitleButton(&data, view), newAlbumDateButton(&data, view)),
			nil,
			container.NewGridWithRows(1, newAlbumAdderLocalButton(&data, view), newAlbumAdderOnlineButton(&data, view)),
			newAlbumSearchBar(&data, view),
		),
		nil,
		nil,
		nil,
		view,
	)
	return container.NewTabItemWithIcon(albumTabTitle, resource.AlbumTabIcon(), border)
}

func newAlbumViewList(data *cbinding.DataList[player.Album]) *cwidget.AlbumViewList {
	albumViewSize := fyne.NewSize(135.0, 165.0)
	list := cwidget.NewAlbumViewList(func(album *player.Album) fyne.CanvasObject {
		view := cwidget.NewAlbumView(album)

		view.OnTapped = func(*fyne.PointEvent) {
			manager.GetCurrentAlbum().Set(album)
		}

		view.OnTappedSecondary = func(event *fyne.PointEvent) {
			newAlbumMenu(fyne.CurrentApp().Driver().CanvasForObject(view), album).ShowAtPosition(event.AbsolutePosition)
		}

		return view
	}, albumViewSize)

	data.Attach(list)
	return list
}

func newAlbumSearchBar(data *cbinding.DataList[player.Album], view *cwidget.AlbumViewList) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(a player.Album) bool {
			return strings.Contains(strings.ToLower(a.Title), title)
		})
	}
	return entry
}

func newAlbumAdderLocalButton(data *cbinding.DataList[player.Album], view *cwidget.AlbumViewList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.AlbumAdderLocalIcon(), func() { showErrorIfAny(manager.AddAlbum()) })
	button.Importance = widget.LowImportance
	return button
}

func newAlbumAdderOnlineButton(data *cbinding.DataList[player.Album], view *cwidget.AlbumViewList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.AlbumAdderOnlineIcon(), func() {})
	button.Importance = widget.LowImportance
	return button
}

func newAlbumTitleButton(data *cbinding.DataList[player.Album], view *cwidget.AlbumViewList) *widget.Button {
	reverse := false
	button := widget.NewButton("Title", func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Album) bool {
			return (strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) < 0) != reverse
		})
	})
	button.Importance = widget.LowImportance
	return button
}

func newAlbumDateButton(data *cbinding.DataList[player.Album], view *cwidget.AlbumViewList) *widget.Button {
	reverse := true
	button := widget.NewButton("Date", func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Album) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
	button.Importance = widget.LowImportance
	button.OnTapped()
	return button
}

func newAlbumMenu(canvas fyne.Canvas, album *player.Album) *widget.PopUpMenu {
	rename := fyne.NewMenuItem("Rename", makeRenameDialog(album))
	cover := fyne.NewMenuItem("Cover", makeCoverDialog(album))
	delete := fyne.NewMenuItem("Delete", makeDeleteAlbumDialog(album))
	return widget.NewPopUpMenu(fyne.NewMenu("", rename, cover, delete), canvas)
}

func makeRenameDialog(album *player.Album) func() {
	entry := widget.NewEntry()
	return func() {
		dialog.ShowCustomConfirm("Enter title:", "Confirm", "Cancel", entry, func(rename bool) {
			if rename {
				log.Printf("rename %v to %v\n", album.Title, entry.Text)
				showErrorIfAny(manager.UpdateTitle(album, entry.Text))
			}
		}, getMainWindow())
	}
}

func makeCoverDialog(album *player.Album) func() {
	return func() {
		fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorIfAny(err)
			} else if result != nil {
				log.Printf("update %v's cover: %v\n", album.Title, result.URI().Path())
				showErrorIfAny(manager.UpdateCover(album, result.URI().Path()))
			}
		}, getMainWindow())
		fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
		fileOpenDialog.SetConfirmText("Upload")
		fileOpenDialog.Show()
	}
}

func makeDeleteAlbumDialog(album *player.Album) func() {
	return func() {
		dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", album.Title), func(delete bool) {
			if delete {
				log.Printf("delete %v\n", album.Title)
				showErrorIfAny(manager.DeleteAlbum(album))
			}
		}, getMainWindow())
	}
}
