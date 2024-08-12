package view

import (
	"fmt"
	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view/internal/cwidget"
	"meowyplayer/view/internal/resource"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicPage struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[[]model.Music]

	pipeline dataPipeline[model.Music]
	current  model.Album
}

func newMusicPage() *MusicPage {
	var p MusicPage
	p = MusicPage{
		searchBar: cwidget.NewSearchBar(
			cwidget.NewCachedList(cwidget.NewMusicCardConstructor(p.playMusic, p.showMusicMenu)),
			p.setEntryFilter,
			nil,
		),
		pipeline: newDataPipeline[model.Music](),
	}

	//search bar menu and toolbar
	p.searchBar.AddDropDown(cwidget.NewMenuItem(resource.MostRecentText(), theme.HistoryIcon(), p.setDateComparator))
	p.searchBar.AddDropDown(cwidget.NewMenuItem(resource.AlphabeticalText(), resource.AlphabeticalIcon(), p.setTitleComparator))
	p.searchBar.AddToolbar(cwidget.NewButton(resource.BackText(), theme.NavigateBackIcon(), model.UIClient().FocusAlbumView))

	//client update callback
	model.UIClient().OnAlbumSelected().Attach(&p)                                         //update current album and list content when selecting album
	model.UIClient().OnStorageLoaded().AttachFunc(func([]model.Album) { p.updateList() }) //update list content when albums get updated

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *MusicPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.searchBar)
}

func (p *MusicPage) Notify(album model.Album) {
	if p.current.Key() != album.Key() {
		p.current = album
		p.searchBar.ClearSearchEntry()
	}
}

func (p *MusicPage) updateList() {
	if p.current.Key().IsEmpty() {
		return
	}

	var err error
	p.current, err = model.UIClient().GetAlbum(p.current.Key())
	if err != nil {
		fyne.LogError("musicPage updateList fails", err)
	}
	p.searchBar.Update(p.pipeline.apply(p.current.Music()))
}

func (p *MusicPage) setEntryFilter(substr string) {
	p.pipeline.filter = func(str string) bool {
		return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
	}
	p.updateList()
}

func (p *MusicPage) setDateComparator() {
	p.pipeline.comparator = func(l, r model.Music) int {
		return -l.Date().Compare(r.Date())
	}
	p.updateList()
}

func (p *MusicPage) setTitleComparator() {
	p.pipeline.comparator = func(l, r model.Music) int {
		return strings.Compare(strings.ToLower(l.Title()), strings.ToLower(r.Title()))
	}
	p.updateList()
}

func (p *MusicPage) playMusic(music model.Music) {
	playlist := p.pipeline.sortCopy(p.current.Music()) //p.current might be outdated
	toPlay := slices.Index(playlist, music)
	player.Instance().LoadAlbum(p.current.Key(), playlist, toPlay)
}

func (p *MusicPage) showMusicMenu(e *fyne.PointEvent, music model.Music) {
	deleteMenu := cwidget.NewMenuItem(resource.DeleteText(), theme.DeleteIcon(), func() { p.showDeleteMusicDialog(music) })
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
}

func (p *MusicPage) showDeleteMusicDialog(music model.Music) {
	dialog.ShowCustomConfirm(resource.DeleteConfirmationText(), resource.DeleteText(), resource.CancelText(),
		widget.NewLabel(fmt.Sprintf(resource.DeleteMusicTextTemplate(), music.Title())),
		func(confirm bool) {
			if confirm {
				if err := model.UIClient().RemoveMusicFromAlbum(p.current.Key(), music.Key()); err != nil {
					fyne.LogError("failed to remove music", err)
				}
			}
		},
		getWindow(),
	)
}
