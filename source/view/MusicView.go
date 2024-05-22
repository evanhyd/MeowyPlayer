package view

import (
	"fmt"
	"playground/cwidget"
	"playground/model"
	"playground/resource"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicView struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[model.Music]
	cover     *canvas.Image
	title     widget.TextSegment
	info      widget.TextSegment
	list      *cwidget.ScrollList[model.Music]

	client  *model.MusicClient
	current model.Album
}

func NewMusicView(client *model.MusicClient) *MusicView {
	var v MusicView
	v = MusicView{
		searchBar: cwidget.NewSearchBar[model.Music](
			v.render,
			cwidget.NewButtonWithIcon(resource.KReturnText, theme.NavigateBackIcon(), client.FocusAlbumView),
		),
		cover: canvas.NewImageFromResource(theme.BrokenImageIcon()),
		title: widget.TextSegment{Style: widget.RichTextStyleHeading},
		info:  widget.TextSegment{Style: widget.RichTextStyleParagraph},
		list: cwidget.NewScrollList(
			container.NewVBox(),
			func() cwidget.WidgetObserver[model.Music] { return newMusicCard() },
		),
		client: client,
	}

	//search bar
	v.searchBar.AddComparator(resource.KMostRecentText, theme.HistoryIcon(), func(a, b model.Music) int {
		return -a.Date().Compare(b.Date())
	})
	v.searchBar.AddComparator(resource.KAlphabeticalText, resource.AlphabeticalIcon, func(a, b model.Music) int {
		return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title()))
	})
	v.searchBar.Select(0)

	//cover
	v.cover.SetMinSize(resource.KAlbumCoverSize)

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
	desc := widget.NewRichText(&v.title, &v.info)
	desc.Wrapping = fyne.TextWrapWord
	desc.Truncation = fyne.TextTruncateEllipsis

	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewBorder(v.searchBar, nil, v.cover, nil, desc), nil, nil, nil, v.list,
	))
}

func (v *MusicView) showDeleteMusicDialog(music model.Music) {
	dialog.ShowCustomConfirm(resource.KDeleteConfirmationText, resource.KDeleteText, resource.KCancelText,
		widget.NewLabel(fmt.Sprintf(resource.KDeleteMusicTextTemplate, music.Title())),
		func(confirm bool) {
			if confirm {
				if err := v.client.RemoveMusic(v.current.Key(), music.Key()); err != nil {
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

	//update album cover and description
	v.cover.Resource = v.current.Cover()
	v.title.Text = v.current.Title()
	v.info.Text = fmt.Sprintf(resource.KAlbumTipTextTemplate, v.current.Count(), v.current.Date().Format(time.DateTime))

	//update music cards
	v.list.Notify(v.searchBar.Query(v.current.Music()))
}
