package view

import (
	"fmt"
	"playground/model"
	"playground/resource"
	"playground/view/internal/cwidget"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AlbumView struct {
	widget.BaseWidget
	list *cwidget.SearchList[model.Album, *AlbumCard]

	client     *model.Client
	pipeline   DataPipeline[model.Album]
	albumsData []model.Album
}

func NewAlbumView(client *model.Client) *AlbumView {
	var v *AlbumView
	v = &AlbumView{
		list: cwidget.NewSearchList(
			container.NewGridWrap(resource.KAlbumCardSize),
			newAlbumCard,
			func(e cwidget.ItemTapEvent[model.Album]) {
				v.client.SelectAlbum(e.Data)
			},
			func(e cwidget.ItemTapEvent[model.Album]) {
				editMenu := cwidget.NewMenuItemWithIcon(resource.KEditText, theme.DocumentCreateIcon(), func() { v.showEditAlbumDialog(e.Data) })
				deleteMenu := cwidget.NewMenuItemWithIcon(resource.KDeleteText, theme.DeleteIcon(), func() { v.showDeleteAlbumDialog(e.Data) })
				widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
			},
			func(sub string) {
				v.pipeline.filter = func(str string) bool {
					return strings.Contains(strings.ToLower(str), strings.ToLower(sub))
				}
				v.updateList()
			},
			nil,
		),
		client: client,
		pipeline: DataPipeline[model.Album]{
			comparator: func(_, _ model.Album) int { return -1 },
			filter:     func(_ string) bool { return true },
		},
	}

	v.list.AddDropDown(cwidget.NewMenuItemWithIcon(resource.KMostRecentText, theme.HistoryIcon(), func() {
		v.pipeline.comparator = func(l, r model.Album) int {
			return -l.Date().Compare(r.Date())
		}
		v.updateList()
	}))
	v.list.AddDropDown(cwidget.NewMenuItemWithIcon(resource.KAlphabeticalText, resource.AlphabeticalIcon, func() {
		v.pipeline.comparator = func(l, r model.Album) int {
			return strings.Compare(strings.ToLower(l.Title()), strings.ToLower(r.Title()))
		}
		v.updateList()
	}))
	v.list.AddToolbar(widget.NewToolbarAction(theme.FolderNewIcon(), v.showCreateAlbumDialog))
	v.ExtendBaseWidget(v)

	client.OnAlbumsChanged().Attach(v)
	return v
}

func (v *AlbumView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.list)
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
	v.albumsData = albums
	v.updateList()
}

func (v *AlbumView) updateList() {
	v.list.Update(v.pipeline.pass(v.albumsData))
}
