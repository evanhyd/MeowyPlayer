package ui

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/utility"
)

func newAlbumTab() *container.TabItem {
	const (
		albumTabTitle    = "Album"
		albumTabIconName = "album_tab.png"
	)

	//album views
	data := cbinding.NewAlbumList()
	view := newAlbumView(data)
	manager.GetCurrentConfig().Attach(data)

	searchBar := newAlbumSearchBar(data, view)
	albumAdderLocalButton := newAlbumAdderLocalButton(data, view)
	albumAdderOnlineButton := newAlbumAdderOnlineButton(data, view)
	titleButton := newAlbumTitleButton(data, view)
	dateButton := newAlbumDateButton(data, view)
	dateButton.OnTapped()

	border := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, titleButton, dateButton),
			nil,
			container.NewGridWithRows(1, albumAdderLocalButton, albumAdderOnlineButton),
			searchBar,
		),
		nil,
		nil,
		nil,
		view,
	)
	return container.NewTabItemWithIcon(albumTabTitle, resource.GetAsset(albumTabIconName), border)
}

func newAlbumView(data binding.DataList) *widget.List {
	const albumCoverIconName = "default.png"
	albumCoverIconSize := fyne.NewSize(128.0, 128.0)
	albumCoverIcon := resource.GetAsset(albumCoverIconName)

	view := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			setting := widget.NewButton(">", func() {})
			setting.Importance = widget.LowImportance
			cover := widget.NewCard("", "", nil)
			cover.Image = canvas.NewImageFromResource(albumCoverIcon)
			intro := widget.NewLabel("")
			return container.NewBorder(nil, nil, container.NewBorder(nil, nil, setting, nil, cover), nil, intro)
		},
		func(item binding.DataItem, canvasObject fyne.CanvasObject) {
			data, err := item.(binding.Untyped).Get()
			utility.MustOk(err)
			album := data.(player.Album)

			objects0 := canvasObject.(*fyne.Container).Objects
			intro := objects0[0].(*widget.Label)

			//optionally update
			if description := album.Description(); intro.Text != description {
				intro.Text = description

				objects1 := objects0[1].(*fyne.Container).Objects

				//update album cover
				cover := objects1[0].(*widget.Card)
				utility.MustNotNil(cover)
				cover.Image = album.Cover
				cover.Image.SetMinSize(albumCoverIconSize)

				//update setting functionality
				setting := objects1[1].(*widget.Button)
				utility.MustNotNil(setting)
				setting.OnTapped = func() {
					newAlbumMenu(fyne.CurrentApp().Driver().CanvasForObject(canvasObject), &album).ShowAtPosition(
						fyne.CurrentApp().Driver().AbsolutePositionForObject(setting))
				}

				canvasObject.Refresh()
			}
		})

	//select and load album
	view.OnSelected = func(id widget.ListItemID) {
		item, err := data.GetItem(id)
		utility.MustOk(err)
		data, err := item.(binding.Untyped).Get()
		utility.MustOk(err)
		a := data.(player.Album)
		manager.GetCurrentAlbum().Set(&a)
		view.Unselect(id)
	}

	return view
}

func newAlbumSearchBar(data *cbinding.AlbumList, view *widget.List) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(title string) {
		title = strings.ToLower(title)
		filter := func(a player.Album) bool {
			return strings.Contains(strings.ToLower(a.Title), title)
		}
		data.SetFilter(filter)
		view.ScrollToTop()
	}
	return entry
}

func newAlbumAdderLocalButton(data *cbinding.AlbumList, view *widget.List) *widget.Button {
	const iconName = "album_adder_local.png"
	button := widget.NewButtonWithIcon("", resource.GetAsset(iconName), func() { showErrorIfAny(manager.AddAlbum()) })
	button.Importance = widget.LowImportance
	return button
}

func newAlbumAdderOnlineButton(data *cbinding.AlbumList, view *widget.List) *widget.Button {
	const iconName = "album_adder_online.png"
	button := widget.NewButtonWithIcon("", resource.GetAsset(iconName), func() {})
	button.Importance = widget.LowImportance
	return button
}

func newAlbumTitleButton(data *cbinding.AlbumList, view *widget.List) *widget.Button {
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

func newAlbumDateButton(data *cbinding.AlbumList, view *widget.List) *widget.Button {
	reverse := true
	button := widget.NewButton("Date", func() {
		reverse = !reverse
		data.SetSorter(func(a1, a2 player.Album) bool {
			return a1.Date.After(a2.Date) != reverse
		})
	})
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
