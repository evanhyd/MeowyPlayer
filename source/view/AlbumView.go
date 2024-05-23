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
	list      *cwidget.ScrollList[model.Album]

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
		list: cwidget.NewScrollList(
			container.NewGridWrap(resource.KAlbumCardSize),
			func() cwidget.WidgetObserver[model.Album] { return newAlbumCard() },
		),
		client: client,
	}

	//search bar
	v.searchBar.AddMenuItem(resource.KMostRecentText, theme.HistoryIcon(), func(a, b model.Album) int {
		return -a.Date().Compare(b.Date())
	})
	v.searchBar.AddMenuItem(resource.KAlphabeticalText, resource.AlphabeticalIcon, func(a, b model.Album) int {
		return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title()))
	})
	v.searchBar.Select(0)

	//item list
	v.list.OnItemTapped = func(e cwidget.ItemTapEvent[model.Album]) {
		v.client.SelectAlbum(e.Data)
	}
	v.list.OnItemTappedSecondary = func(e cwidget.ItemTapEvent[model.Album]) {
		editMenu := cwidget.NewMenuItemWithIcon(resource.KEditText, theme.DocumentCreateIcon(), func() { v.showEditAlbumDialog(e.Data) })
		deleteMenu := cwidget.NewMenuItemWithIcon(resource.KDeleteText, theme.DeleteIcon(), func() { v.showDeleteAlbumDialog(e.Data) })
		widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
	}

	client.OnAlbumsChanged().Attach(&v)
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *AlbumView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(v.searchBar, nil, nil, nil, v.list))
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
	v.list.Notify(v.searchBar.Query(v.albums))
}
