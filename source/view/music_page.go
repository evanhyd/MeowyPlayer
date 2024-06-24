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

type MusicPage struct {
	widget.BaseWidget
	list *cwidget.SearchList[model.Music, *cwidget.MusicCard]

	pipeline DataPipeline[model.Music]
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
		pipeline: NewDataPipeline[model.Music](),
	}

	//search bar menu and toolbar
	p.list.AddDropDown(cwidget.NewMenuItem(resource.KMostRecentText, theme.HistoryIcon(), p.setDateComparator))
	p.list.AddDropDown(cwidget.NewMenuItem(resource.KAlphabeticalText, resource.AlphabeticalIcon, p.setTitleComparator))
	p.list.AddToolbar(cwidget.NewButton(resource.KBackText, theme.NavigateBackIcon(), model.GetClient().FocusAlbumView))

	//client update callback
	model.GetClient().OnAlbumSelected().Attach(&p)                                         //update current album and list content when selecting album
	model.GetClient().OnAlbumsChanged().AttachFunc(func([]model.Album) { p.updateList() }) //update list content when albums get updated

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
	p.current, err = model.GetClient().GetAlbum(p.current.Key())
	if err != nil {
		fyne.LogError("musicPage updateList fails", err)
	}
	p.list.Update(p.pipeline.Pass(p.current.Music()))
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
	fmt.Println(p.current, music)
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
				if err := model.GetClient().RemoveMusicFromAlbum(p.current.Key(), music.Key()); err != nil {
					fyne.LogError("failed to remove music", err)
				}
			}
		},
		getWindow(),
	)
}
