package view

import (
	"fmt"
	"playground/cwidget"
	"playground/model"
	"playground/resource"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AlbumView struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[model.Album]
	cards     *cwidget.ScrollList[AlbumCardProp]

	client *model.MusicClient
	albums []model.Album
}

func NewAlbumView(client *model.MusicClient) *AlbumView {
	var v AlbumView
	v = AlbumView{
		searchBar: cwidget.NewSearchBar[model.Album](
			v.render,
			cwidget.NewButtonWithIcon("", theme.FolderNewIcon(), v.showCreateAlbumDialog),
		),
		cards: cwidget.NewScrollList(
			container.NewGridWrap(resource.KAlbumCardSize),
			func() cwidget.ObserverCanvasObject[AlbumCardProp] { return newAlbumCard() },
		),
		client: client,
	}

	v.searchBar.AddComparator(resource.KMostRecentMenuText, theme.HistoryIcon(), func(a, b model.Album) int {
		return -a.Date().Compare(b.Date())
	})
	v.searchBar.AddComparator(resource.KAlphabeticalMenuText, resource.AlphabeticalIcon, func(a, b model.Album) int {
		return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title()))
	})
	v.searchBar.Select(0)

	client.OnAlbumsChanged().Attach(&v)
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *AlbumView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(v.searchBar, nil, nil, nil, v.cards))
}

func (v *AlbumView) showCreateAlbumDialog() {
	editor := newAlbumEditor()
	dialog.ShowCustomConfirm(resource.KCreateAlbumText, resource.KCreateText, resource.KCancelText, editor,
		func(confirm bool) {
			if confirm {
				if err := v.client.CreateAlbum(editor.state()); err != nil {
					fyne.LogError("failed to create album", err)
				}
			}
		},
		getWindow(),
	)
}

func (v *AlbumView) showEditAlbumDialog(album model.Album) {
	editor := newAlbumEditorWithState(album.Title(), album.Cover())
	dialog.ShowCustomConfirm(resource.KEditAlbumText, resource.KSaveText, resource.KCancelText, editor,
		func(confirm bool) {
			if confirm {
				title, cover := editor.state()
				if err := v.client.EditAlbum(album.Key(), title, cover); err != nil {
					fyne.LogError("failed to edit album", err)
				}
			}
		},
		getWindow(),
	)
}

func (v *AlbumView) showDeleteAlbumDialog(album model.Album) {
	dialog.ShowCustomConfirm(resource.KDeleteConfirmationText, resource.KDeleteText, resource.KCancelText,
		widget.NewLabel(fmt.Sprintf(resource.KDeleteAlbumTextTemplate, album.Title())),
		func(confirm bool) {
			if confirm {
				if err := v.client.RemoveAlbum(album.Key()); err != nil {
					fyne.LogError("failed to remove album", err)
				}
			}
		},
		getWindow(),
	)
}

func (v *AlbumView) Notify(albums []model.Album) {
	v.albums = albums
	v.render()
}

func (v *AlbumView) render() {
	albums := v.searchBar.Query(v.albums)

	//update collection view
	props := make([]AlbumCardProp, 0, len(albums))
	for _, album := range albums {
		album := album //loop closure

		editMenu := fyne.NewMenuItem(resource.KEditText, func() { v.showEditAlbumDialog(album) })
		editMenu.Icon = theme.DocumentCreateIcon()
		deleteMenu := fyne.NewMenuItem(resource.KDeleteText, func() { v.showDeleteAlbumDialog(album) })
		deleteMenu.Icon = theme.DeleteIcon()

		props = append(props, AlbumCardProp{
			Album:    album,
			OnTapped: func(*fyne.PointEvent) { v.client.FocusAlbum(album) },
			OnTappedSecondary: func(e *fyne.PointEvent) {
				widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
			},
		})
	}
	v.cards.Notify(props)
}
