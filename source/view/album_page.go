package view

import (
	"fmt"
	"meowyplayer/model"
	"meowyplayer/view/internal/cwidget"
	"meowyplayer/view/internal/resource"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AlbumPage struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[[]model.Album]

	pipeline   dataPipeline[model.Album]
	albumsData []model.Album
}

func newAlbumPage() *AlbumPage {
	var p AlbumPage
	p = AlbumPage{
		searchBar: cwidget.NewSearchBar(
			cwidget.NewCustomList(container.NewGridWrap(resource.KAlbumCardSize), cwidget.NewAlbumCardConstructor(p.selectAlbum, p.showAlbumMenu)),
			p.setEntryFilter,
			nil,
		),
		pipeline: newDataPipeline[model.Album](),
	}

	//search bar menu and toolbar
	p.searchBar.AddDropDown(cwidget.NewMenuItem(resource.MostRecentText(), theme.HistoryIcon(), p.setDateComparator))
	p.searchBar.AddDropDown(cwidget.NewMenuItem(resource.AlphabeticalText(), resource.AlphabeticalIcon(), p.setTitleComparator))
	p.searchBar.AddToolbar(cwidget.NewButton(resource.CreateAlbumText(), theme.FolderNewIcon(), p.showCreateAlbumDialog))
	p.ExtendBaseWidget(&p)

	//client update callback
	model.Instance().OnStorageLoaded().Attach(&p)
	return &p
}

func (p *AlbumPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.searchBar)
}

func (p *AlbumPage) Notify(albums []model.Album) {
	p.albumsData = albums
	p.updateList()
}

func (p *AlbumPage) updateList() {
	p.searchBar.Update(p.pipeline.apply(p.albumsData))
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
	err := model.Instance().SelectAlbum(key)
	if err != nil {
		fyne.LogError("selectAlbum fail", err)
	}
}

func (p *AlbumPage) showAlbumMenu(e *fyne.PointEvent, key model.AlbumKey) {
	album, err := model.Instance().GetAlbum(key)
	if err != nil {
		fyne.LogError("showAlbumMenu fail", err)
	}

	editMenu := cwidget.NewMenuItem(resource.EditText(), theme.DocumentCreateIcon(), func() { p.showEditAlbumDialog(album) })
	deleteMenu := cwidget.NewMenuItem(resource.DeleteText(), theme.DeleteIcon(), func() { p.showDeleteAlbumDialog(album) })
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
}

func (p *AlbumPage) showEditAlbumDialog(album model.Album) {
	editor := newAlbumEditorWithState(album.Title(), album.Cover())
	dialog.ShowCustomConfirm(resource.EditAlbumText(), resource.SaveText(), resource.CancelText(), editor,
		func(confirm bool) {
			if confirm {
				title, cover := editor.state()
				if err := model.Instance().EditAlbum(album.Key(), title, cover); err != nil {
					fyne.LogError("failed to edit album", err)
				}
			}
		},
		getWindow(),
	)
}

func (p *AlbumPage) showDeleteAlbumDialog(album model.Album) {
	dialog.ShowCustomConfirm(resource.DeleteConfirmationText(), resource.DeleteText(), resource.CancelText(),
		widget.NewLabel(fmt.Sprintf(resource.DeleteAlbumTextTemplate(), album.Title())),
		func(confirm bool) {
			if confirm {
				if err := model.Instance().RemoveAlbum(album.Key()); err != nil {
					fyne.LogError("failed to remove album", err)
				}
			}
		},
		getWindow(),
	)
}

func (p *AlbumPage) showCreateAlbumDialog() {
	editor := newAlbumEditor()
	dialog.ShowCustomConfirm(resource.CreateAlbumText(), resource.CreateText(), resource.CancelText(), editor,
		func(confirm bool) {
			if confirm {
				if err := model.Instance().CreateAlbum(editor.state()); err != nil {
					fyne.LogError("failed to create album", err)
				}
			}
		},
		getWindow(),
	)
}
