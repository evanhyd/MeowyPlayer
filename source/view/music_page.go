package view

import (
	"fmt"
	"playground/model"
	"playground/player"
	"playground/view/internal/cwidget"
	"playground/view/internal/resource"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicPage struct {
	widget.BaseWidget
	list *cwidget.SearchList[model.Music, *cwidget.MusicCard]

	pipeline dataPipeline[model.Music]
	current  model.Album
}

func newMusicPage() *MusicPage {
	var p MusicPage
	p = MusicPage{
		list: cwidget.NewSearchList(
			container.NewVBox(),
			cwidget.NewMusicCardConstructor(p.playMusic, p.showMusicMenu),
			p.setEntryFilter,
			nil,
		),
		pipeline: newDataPipeline[model.Music](),
	}

	//search bar menu and toolbar
	p.list.AddDropDown(cwidget.NewMenuItem(resource.KMostRecentText, theme.HistoryIcon(), p.setDateComparator))
	p.list.AddDropDown(cwidget.NewMenuItem(resource.KAlphabeticalText, resource.AlphabeticalIcon, p.setTitleComparator))
	p.list.AddToolbar(cwidget.NewButton(resource.KBackText, theme.NavigateBackIcon(), model.Instance().FocusAlbumView))

	//client update callback
	model.Instance().OnAlbumSelected().Attach(&p)                                         //update current album and list content when selecting album
	model.Instance().OnAlbumsChanged().AttachFunc(func([]model.Album) { p.updateList() }) //update list content when albums get updated

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *MusicPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.list)
}

func (p *MusicPage) Notify(album model.Album) {
	if p.current.Key() != album.Key() {
		p.current = album
		p.list.ClearSearchEntry()
	}
}

func (p *MusicPage) updateList() {
	if p.current.Key().IsEmpty() {
		return
	}

	var err error
	p.current, err = model.Instance().GetAlbum(p.current.Key())
	if err != nil {
		fyne.LogError("musicPage updateList fails", err)
	}
	p.list.Update(p.pipeline.apply(p.current.Music()))
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
	playlist := p.pipeline.sortCopy(p.current.Music())
	toPlay := slices.Index(playlist, music)
	player.Instance().LoadAlbum(playlist, toPlay)
}

func (p *MusicPage) showMusicMenu(e *fyne.PointEvent, music model.Music) {
	deleteMenu := cwidget.NewMenuItem(resource.KDeleteText, theme.DeleteIcon(), func() { p.showDeleteMusicDialog(music) })
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
}

func (p *MusicPage) showDeleteMusicDialog(music model.Music) {
	dialog.ShowCustomConfirm(resource.KDeleteConfirmationText, resource.KDeleteText, resource.KCancelText,
		widget.NewLabel(fmt.Sprintf(resource.KDeleteMusicTextTemplate, music.Title())),
		func(confirm bool) {
			if confirm {
				if err := model.Instance().RemoveMusicFromAlbum(p.current.Key(), music.Key()); err != nil {
					fyne.LogError("failed to remove music", err)
				}
			}
		},
		getWindow(),
	)
}
