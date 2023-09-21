package ui

import (
	"fmt"
	"log"

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
)

func newAlbumTab() *container.TabItem {
	data := cbinding.MakeAlbumDataList()
	manager.GetCurrentConfig().Attach(&data)

	return container.NewTabItemWithIcon("Album", resource.AlbumTabIcon(), container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, cwidget.NewAlbumTitleButton(&data, "Title"), cwidget.NewAlbumDateButton(&data, "Date")),
			nil,
			container.NewGridWithRows(1, newAlbumAdderLocalButton(&data), newAlbumAdderOnlineButton(&data)),
			cwidget.NewAlbumSearchBar(&data),
		),
		nil,
		nil,
		nil,
		newAlbumViewList(&data),
	))
}

func newAlbumViewList(data *cbinding.AlbumDataList) *cwidget.AlbumViewList {
	return cwidget.NewAlbumViewList(data, func(album *player.Album) fyne.CanvasObject {
		view := cwidget.NewAlbumView(album)
		view.OnTapped = func(*fyne.PointEvent) { manager.GetCurrentAlbum().Set(album) }
		view.OnTappedSecondary = func(event *fyne.PointEvent) {
			canvas := fyne.CurrentApp().Driver().CanvasForObject(view)
			newAlbumMenu(canvas, album).ShowAtPosition(event.AbsolutePosition)
		}
		return view
	}, fyne.NewSize(135.0, 165.0))
}

func newAlbumAdderLocalButton(data *cbinding.AlbumDataList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.AlbumAdderLocalIcon(), func() { showErrorIfAny(manager.AddAlbum()) })
	button.Importance = widget.LowImportance
	return button
}

func newAlbumAdderOnlineButton(data *cbinding.AlbumDataList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.AlbumAdderOnlineIcon(), func() {})
	button.Importance = widget.LowImportance
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
