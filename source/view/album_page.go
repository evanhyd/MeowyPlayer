package view

import (
	"fmt"
	"playground/model"
	"playground/view/internal/cwidget"
	"playground/view/internal/resource"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AlbumPage struct {
	widget.BaseWidget
	list *cwidget.SearchList[model.Album, *cwidget.AlbumCard]

	pipeline   DataPipeline[model.Album]
	albumsData []model.Album
}

func newAlbumPage() *AlbumPage {
	var p AlbumPage
	p = AlbumPage{
		list: cwidget.NewSearchList(
			container.NewGridWrap(resource.KAlbumCardSize),
			cwidget.NewAlbumCardConstructor(p.selectAlbum, p.showAlbumMenu),
			p.setEntryFilter,
			nil,
		),
		pipeline: newDataPipeline[model.Album](),
	}

	//search bar menu and toolbar
	p.list.AddDropDown(cwidget.NewMenuItem(resource.KMostRecentText, theme.HistoryIcon(), p.setDateComparator))
	p.list.AddDropDown(cwidget.NewMenuItem(resource.KAlphabeticalText, resource.AlphabeticalIcon, p.setTitleComparator))
	p.list.AddToolbar(cwidget.NewButton(resource.KCreateAlbumText, theme.FolderNewIcon(), p.showCreateAlbumDialog))
	p.ExtendBaseWidget(&p)

	//client update callback
	model.GetClient().OnAlbumsChanged().Attach(&p)
	return &p
}

func (p *AlbumPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.list)
}

func (p *AlbumPage) Notify(albums []model.Album) {
	p.albumsData = albums
	p.updateList()
}

func (p *AlbumPage) updateList() {
	p.list.Update(p.pipeline.pass(p.albumsData))
}

func (p *AlbumPage) setEntryFilter(substr string) {
	p.pipeline.filter = func(str string) bool {
		return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
	}
	p.updateList()
}

func (p *AlbumPage) setDateComparator() {
	p.pipeline.comparator = func(l, r model.Album) int {
		return -l.Date().Compare(r.Date())
	}
	p.updateList()
}

func (p *AlbumPage) setTitleComparator() {
	p.pipeline.comparator = func(l, r model.Album) int {
		return strings.Compare(strings.ToLower(l.Title()), strings.ToLower(r.Title()))
	}
	p.updateList()
}

func (p *AlbumPage) selectAlbum(key model.AlbumKey) {
	err := model.GetClient().SelectAlbum(key)
	if err != nil {
		fyne.LogError("selectAlbum fail", err)
	}
}

func (p *AlbumPage) showAlbumMenu(e *fyne.PointEvent, key model.AlbumKey) {
	album, err := model.GetClient().GetAlbum(key)
	if err != nil {
		fyne.LogError("showAlbumMenu fail", err)
	}

	editMenu := cwidget.NewMenuItem(resource.KEditText, theme.DocumentCreateIcon(), func() { p.showEditAlbumDialog(album) })
	deleteMenu := cwidget.NewMenuItem(resource.KDeleteText, theme.DeleteIcon(), func() { p.showDeleteAlbumDialog(album) })
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
}

func (p *AlbumPage) showEditAlbumDialog(album model.Album) {
	editor := newAlbumEditorWithState(album.Title(), album.Cover())
	dialog.ShowCustomConfirm(resource.KEditAlbumText, resource.KSaveText, resource.KCancelText, editor,
		func(confirm bool) {
			if confirm {
				title, cover := editor.state()
				if err := model.GetClient().EditAlbum(album.Key(), title, cover); err != nil {
					fyne.LogError("failed to edit album", err)
				}
			}
		},
		getWindow(),
	)
}

func (p *AlbumPage) showDeleteAlbumDialog(album model.Album) {
	dialog.ShowCustomConfirm(resource.KDeleteConfirmationText, resource.KDeleteText, resource.KCancelText,
		widget.NewLabel(fmt.Sprintf(resource.KDeleteAlbumTextTemplate, album.Title())),
		func(confirm bool) {
			if confirm {
				if err := model.GetClient().RemoveAlbum(album.Key()); err != nil {
					fyne.LogError("failed to remove album", err)
				}
			}
		},
		getWindow(),
	)
}

func (p *AlbumPage) showCreateAlbumDialog() {
	editor := newAlbumEditor()
	dialog.ShowCustomConfirm(resource.KCreateAlbumText, resource.KCreateText, resource.KCancelText, editor,
		func(confirm bool) {
			if confirm {
				if err := model.GetClient().CreateAlbum(editor.state()); err != nil {
					fyne.LogError("failed to create album", err)
				}
			}
		},
		getWindow(),
	)
}
