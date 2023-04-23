package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

const (
	albumTabName = "Album"
)

var defaultIcon fyne.Resource
var albumTabIcon fyne.Resource
var albumAdderTabIcon fyne.Resource

func init() {
	const (
		albumCoverIconName = "default.png"
		albumTabIconName   = "album_tab.png"
		albumAdderIconName = "album_adder.png"
	)

	var err error
	if defaultIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumCoverIconName)); err != nil {
		log.Fatal(err)
	}

	if albumTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumTabIconName)); err != nil {
		log.Fatal(err)
	}

	if albumAdderTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumAdderIconName)); err != nil {
		log.Fatal(err)
	}
}

func createAblumTab() *container.TabItem {
	albumAdderButton := cwidget.NewButtonWithIcon("", albumAdderTabIcon)
	searchBar := widget.NewEntry()
	sortByTitleButton := cwidget.NewButton("Title")
	sortByDateButton := cwidget.NewButton("Date")

	scroll := cwidget.NewAlbumList(
		func() fyne.CanvasObject {
			card := widget.NewCard("", "", nil)
			card.Image = canvas.NewImageFromResource(defaultIcon)
			card.Image.SetMinSize(resource.GetAlbumCoverSize())
			title := widget.NewLabel("")
			button := cwidget.NewButton("<")
			return container.NewBorder(nil, nil, card, button, title)
		},
		func(album player.Album, canvas fyne.CanvasObject) {

			//not a solid design. If the inner border style change, then this code would break
			items := canvas.(*fyne.Container).Objects
			label := items[0].(*widget.Label)
			if label.Text != album.Description() {
				label.Text = album.Description()

				//update album cover
				card := items[1].(*widget.Card)
				card.Image = album.CoverIcon()

				//update setting menu
				button := items[2].(*cwidget.Button)
				button.OnTapped = func() {
					createAlbumPopUpMenu(fyne.CurrentApp().Driver().CanvasForObject(button), album).
						ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(button))
				}

				canvas.Refresh()
			}
		},
	)

	albumAdderButton.OnTapped = func() { DisplayErrorIfAny(player.AddNewAlbum()) }
	searchBar.OnChanged = scroll.SetTitleFilter
	sortByTitleButton.OnTapped = scroll.SetTitleSorter
	sortByDateButton.OnTapped = scroll.SetDateSorter
	player.GetState().OnUpdateAlbums().AddCallback(scroll.SetItems)
	player.GetState().OnFocusAlbumTab().AddCallback(func() { searchBar.SetText("") })
	scroll.SetOnSelected(func(album *player.Album) { player.UserSelectAlbum(*album) })

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
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
		scroll,
	)
	return container.NewTabItemWithIcon(albumTabName, albumTabIcon, canvas)
}

func createAlbumPopUpMenu(canvas fyne.Canvas, album player.Album) *widget.PopUpMenu {
	rename := fyne.NewMenuItem("Rename", func() {
		entry := widget.NewEntry()
		dialog.ShowCustomConfirm("Enter title:", "Confirm", "Cancel", entry, func(shouldRename bool) {
			if shouldRename {
				DisplayErrorIfAny(player.RenameAlbum(album, entry.Text))
			}
		}, player.GetMainWindow())
	})

	cover := fyne.NewMenuItem("Cover", func() {
		fileOpenDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
			if err != nil {
				DisplayErrorIfAny(err)
				return
			}
			if result != nil {
				DisplayErrorIfAny(player.SetAlbumCover(album, result.URI().Path()))
			}
		}, player.GetMainWindow())
		fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
		fileOpenDialog.SetConfirmText("Upload")
		fileOpenDialog.Show()
	})

	delete := fyne.NewMenuItem("Delete", func() {
		dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", album.Title()), func(shouldDelete bool) {
			if shouldDelete {
				DisplayErrorIfAny(player.RemoveAlbum(album))
			}
		}, player.GetMainWindow())
	})
	return widget.NewPopUpMenu(fyne.NewMenu("", rename, cover, delete), canvas)
}
