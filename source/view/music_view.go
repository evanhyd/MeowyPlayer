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

type MusicView struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[model.Music]
	list      *cwidget.ScrollList[model.Music]
	client    *model.Client
	current   model.Album
}

func NewMusicView(client *model.Client) *MusicView {
	var v MusicView
	v = MusicView{
		searchBar: cwidget.NewSearchBar[model.Music](
			v.render,
			cwidget.NewButtonWithIcon(resource.KReturnText, theme.NavigateBackIcon(), client.FocusAlbumView),
		),
		list: cwidget.NewScrollList(
			container.NewVBox(),
			func() cwidget.WidgetObserver[model.Music] { return newMusicCard() },
		),
		client: client,
	}

	//search bar
	v.searchBar.AddMenuItem(resource.KMostRecentText, theme.HistoryIcon(), func() {
		v.searchBar.SetComparator(func(a, b model.Music) int {
			return -a.Date().Compare(b.Date())
		})
	})
	v.searchBar.AddMenuItem(resource.KAlphabeticalText, resource.AlphabeticalIcon, func() {
		v.searchBar.SetComparator(func(a, b model.Music) int {
			return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title()))
		})
	})
	v.searchBar.Select(0)

	//list
	v.list.OnItemTapped = func(e cwidget.ItemTapEvent[model.Music]) {
		//TODO: play music
	}
	v.list.OnItemTappedSecondary = func(e cwidget.ItemTapEvent[model.Music]) {
		deleteMenu := cwidget.NewMenuItemWithIcon(resource.KDeleteText, theme.DeleteIcon(), func() { v.showDeleteMusicDialog(e.Data) })
		widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
	}

	client.OnAlbumSelected().Attach(&v)                                     //update when selecting album
	client.OnAlbumsChanged().AttachFunc(func([]model.Album) { v.render() }) //update when updating album
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *MusicView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(v.searchBar, nil, nil, nil, v.list))
}

func (v *MusicView) showDeleteMusicDialog(music model.Music) {
	dialog.ShowCustomConfirm(resource.KDeleteConfirmationText, resource.KDeleteText, resource.KCancelText,
		widget.NewLabel(fmt.Sprintf(resource.KDeleteMusicTextTemplate, music.Title())),
		func(confirm bool) {
			if confirm {
				if err := v.client.RemoveMusicFromAlbum(v.current.Key(), music.Key()); err != nil {
					fyne.LogError("failed to remove music", err)
				}
			}
		},
		getWindow(),
	)
}

func (v *MusicView) Notify(album model.Album) {
	if v.current.Key() != album.Key() {
		v.current = album
		v.searchBar.ClearFilter() //call render()
	}
}

func (v *MusicView) render() {
	if v.current.Key().IsEmpty() {
		return
	}
	v.current = v.client.GetAlbum(v.current.Key())

	//update music cards
	v.list.Notify(v.searchBar.Query(v.current.Music()))
}
