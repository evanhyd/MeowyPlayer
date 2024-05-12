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
	cards     *cwidget.ScrollList[MusicCardProp]

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
		cover:  canvas.NewImageFromResource(theme.BrokenImageIcon()),
		title:  widget.TextSegment{Style: widget.RichTextStyleHeading},
		info:   widget.TextSegment{Style: widget.RichTextStyleParagraph},
		cards:  cwidget.NewScrollList(container.NewVBox(), func() cwidget.ObserverCanvasObject[MusicCardProp] { return newMusicCard() }),
		client: client,
	}

	v.searchBar.AddComparator(resource.KMostRecentMenuText, theme.HistoryIcon(), func(a, b model.Music) int {
		return -a.Date().Compare(b.Date())
	})
	v.searchBar.AddComparator(resource.KAlphabeticalMenuText, resource.AlphabeticalIcon, func(a, b model.Music) int {
		return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title()))
	})
	v.searchBar.Select(0)

	v.cover.SetMinSize(resource.KAlbumCoverSize)

	client.OnAlbumSelected().Attach(&v)                                     //update when selecting album
	client.OnAlbumsChanged().AttachFunc(func([]model.Album) { v.render() }) //update when updating album
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *MusicView) CreateRenderer() fyne.WidgetRenderer {
	description := widget.NewRichText(&v.title, &v.info)
	description.Wrapping = fyne.TextWrapWord
	description.Truncation = fyne.TextTruncateEllipsis

	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewBorder(v.searchBar, nil, v.cover, nil, description), nil, nil, nil, v.cards,
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
	musicList := v.searchBar.Query(v.current.Music())

	props := make([]MusicCardProp, 0, len(musicList))
	for _, music := range musicList {
		music := music //loop closure

		// editMenu := fyne.NewMenuItem(resource.KEditText, func() { v.showEditAlbumDialog(album) })
		// editMenu.Icon = theme.DocumentCreateIcon()
		deleteMenu := fyne.NewMenuItem(resource.KDeleteText, func() { v.showDeleteMusicDialog(music) })
		deleteMenu.Icon = theme.DeleteIcon()

		props = append(props, MusicCardProp{
			Music:    music,
			OnTapped: func(*fyne.PointEvent) {},
			OnTappedSecondary: func(e *fyne.PointEvent) {
				widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
			},
		})
	}
	v.cards.Notify(props)
}
