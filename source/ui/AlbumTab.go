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
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource/album"
	"meowyplayer.com/source/resource/config"
	"meowyplayer.com/source/resource/texture"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/utility"
)

func newAlbumTab() *container.TabItem {
	const (
		albumTabTitle      = "Album"
		albumTabIconName   = "album_tab.png"
		albumAdderIconName = "album_adder.png"
	)

	// album views
	data := cbinding.NewConfigList()
	view := newAlbumView(data)
	config.Current().Attach(data)

	//search bar
	searchBar := widget.NewEntry()
	searchBar.OnChanged = func(title string) {
		data.SetFilter(makeFilter(title))
		view.ScrollToTop()
	}

	//add new album
	albumAdderButton := widget.NewButtonWithIcon("", texture.Get(albumAdderIconName), func() {
		log.Println("create new album")
		showErrorIfAny(album.Make())
	})
	albumAdderButton.Importance = widget.LowImportance

	//sort by title button
	reverseTitle := false
	sortByTitleButton := widget.NewButton("Title", func() {
		reverseTitle = !reverseTitle
		data.SetSorter(makeTitleSorter(reverseTitle))
	})
	sortByTitleButton.Importance = widget.LowImportance

	//sort by date button
	reverseDate := true
	sortByDateButton := widget.NewButton("Date", func() {
		reverseDate = !reverseDate
		data.SetSorter(makeDateSorter(reverseDate))
	})
	sortByDateButton.Importance = widget.LowImportance
	sortByDateButton.OnTapped()

	border := container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, sortByTitleButton, sortByDateButton),
			nil,
			albumAdderButton,
			searchBar,
		),
		nil,
		nil,
		nil,
		view,
	)
	return container.NewTabItemWithIcon(albumTabTitle, texture.Get(albumTabIconName), border)
}

func newAlbumView(data binding.DataList) *widget.List {
	const (
		albumCoverIconName = "default.png"
	)
	albumCoverIconSize := fyne.NewSize(128.0, 128.0)
	albumCoverIcon := texture.Get(albumCoverIconName)

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
		album.Set(&a)
		view.Unselect(id)
	}

	return view
}

func makeFilter(title string) func(player.Album) bool {
	title = strings.ToLower(title)
	return func(a player.Album) bool {
		return strings.Contains(strings.ToLower(a.Title), title)
	}
}

func makeTitleSorter(reverse bool) func(player.Album, player.Album) bool {
	return func(a1, a2 player.Album) bool {
		return (strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) < 0) != reverse
	}
}

func makeDateSorter(reverse bool) func(player.Album, player.Album) bool {
	return func(a1, a2 player.Album) bool {
		return a1.Date.After(a2.Date) != reverse
	}
}

func newAlbumMenu(canvas fyne.Canvas, album *player.Album) *widget.PopUpMenu {
	rename := fyne.NewMenuItem("Rename", makeRenameDialog(album))
	cover := fyne.NewMenuItem("Cover", makeCoverDialog(album))
	delete := fyne.NewMenuItem("Delete", makeDeleteDialog(album))
	return widget.NewPopUpMenu(fyne.NewMenu("", rename, cover, delete), canvas)
}

func makeRenameDialog(selectedAlbum *player.Album) func() {
	entry := widget.NewEntry()
	return func() {
		dialog.ShowCustomConfirm("Enter title:", "Confirm", "Cancel", entry, func(rename bool) {
			if rename {
				log.Printf("rename %v to %v\n", selectedAlbum.Title, entry.Text)
				showErrorIfAny(album.UpdateTitle(selectedAlbum, entry.Text))
			}
		}, getMainWindow())
	}
}

func makeCoverDialog(selectedAlbum *player.Album) func() {
	return func() {
		fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorIfAny(err)
			} else if result != nil {
				log.Printf("update %v's cover: %v\n", selectedAlbum.Title, result.URI().Path())
				showErrorIfAny(album.UpdateCover(selectedAlbum, result.URI().Path()))
			}
		}, getMainWindow())
		fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
		fileOpenDialog.SetConfirmText("Upload")
		fileOpenDialog.Show()
	}
}

func makeDeleteDialog(selectedAlbum *player.Album) func() {
	return func() {
		dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", selectedAlbum.Title), func(delete bool) {
			if delete {
				log.Printf("delete %v\n", selectedAlbum.Title)
				showErrorIfAny(album.Delete(selectedAlbum))
			}
		}, getMainWindow())
	}
}
