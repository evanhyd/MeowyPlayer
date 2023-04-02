package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

const (
	albumTabName = "Album"
)

var albumCoverIconSize fyne.Size
var albumCoverIcon *canvas.Image
var albumTabIcon fyne.Resource

func init() {
	const (
		albumCoverIconName = "album_cover.png"
		albumTabIconName   = "album_tab.png"
	)

	albumCoverIconSize = fyne.NewSize(128.0, 128.0)
	albumCoverIcon = canvas.NewImageFromFile(resource.GetResourcePath(albumCoverIconName))
	albumCoverIcon.SetMinSize(albumCoverIconSize)

	var err error
	if albumTabIcon, err = fyne.LoadResourceFromPath(resource.GetResourcePath(albumTabIconName)); err != nil {
		log.Fatal(err)
	}
}

func createAblumTab() *container.TabItem {
	searchBar := cwidget.NewSearchBar()
	sortByTitleButton := cwidget.NewButton("Title")
	sortByDateButton := cwidget.NewButton("Date")

	scroll := cwidget.NewAlbumItemList(
		func() fyne.CanvasObject {
			card := widget.NewCard("", "", nil)
			card.SetImage(albumCoverIcon)
			title := widget.NewLabel("")
			return container.NewBorder(nil, nil, card, nil, title)
		},
		func(album player.Album, canvas fyne.CanvasObject) {
			//not a solid design. If the inner border style change, then this code would break
			label := canvas.(*fyne.Container).Objects[0].(*widget.Label)
			if label.Text != album.Description() {
				label.SetText(album.Description())
			}

			card := canvas.(*fyne.Container).Objects[1].(*widget.Card)
			if card.Image != album.CoverIcon() {
				card.SetImage(album.CoverIcon())
				card.Image.SetMinSize(albumCoverIconSize)
			}
		},
	)

	searchBar.SetOnChanged(scroll.SetTitleFilter)
	sortByTitleButton.SetOnTapped(scroll.SetTitleSorter)
	sortByDateButton.SetOnTapped(scroll.SetDateSorter)
	player.GetState().OnReadAlbumsFromDisk().AddObserver(scroll)
	scroll.SetOnSelected(player.GetState().SetSelectedAlbum)

	defer sortByDateButton.OnTapped()

	canvas := container.NewBorder(
		container.NewBorder(
			searchBar,
			nil,
			nil,
			nil,
			container.NewGridWithRows(1, sortByTitleButton, sortByDateButton),
		),
		nil,
		nil,
		nil,
		scroll,
	)
	return container.NewTabItemWithIcon(albumTabName, albumTabIcon, canvas)
}
