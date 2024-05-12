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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicView struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[model.Music]
	cover     *canvas.Image
	info      *widget.RichText
	cards     *cwidget.ScrollList[MusicCardProp]

	client  *model.MusicClient
	current model.Album
}

func NewMusicView(client *model.MusicClient) *MusicView {
	var v MusicView
	v = MusicView{
		searchBar: cwidget.NewSearchBar[model.Music](
			v.render,
			cwidget.NewButtonWithIcon(resource.KReturnText, theme.NavigateBackIcon(), client.RefreshAlbums),
		),
		cover: canvas.NewImageFromResource(theme.BrokenImageIcon()),
		info: widget.NewRichText(
			&widget.TextSegment{Text: "?", Style: widget.RichTextStyle{SizeName: theme.SizeNameHeadingText}},
			&widget.TextSegment{Text: "?", Style: widget.RichTextStyle{SizeName: theme.SizeNameText}},
		),
		cards: cwidget.NewScrollList(container.NewVBox(),
			func() cwidget.ObserverCanvasObject[MusicCardProp] { return newMusicCard() },
		),
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

	v.info.Wrapping = fyne.TextWrapWord
	v.info.Truncation = fyne.TextTruncateEllipsis

	client.OnAlbumFocused().Attach(&v)                                        //update when switching album
	client.OnAlbumsChanged().AttachFunc(func(_ []model.Album) { v.render() }) //update when updating album
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *MusicView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewBorder(v.searchBar, nil, v.cover, nil, v.info), nil, nil, nil, v.cards,
	))
}

func (v *MusicView) Notify(album model.Album) {
	if v.current.Key() != album.Key() {
		v.current = album
		v.render()
	}
}

func (v *MusicView) render() {
	if v.current.Key().IsEmpty() {
		return
	}
	v.current = v.client.GetAlbum(v.current.Key())

	//update album cover and description
	v.cover.Resource = v.current.Cover()
	v.info.Segments[0].(*widget.TextSegment).Text = v.current.Title()
	v.info.Segments[1].(*widget.TextSegment).Text = fmt.Sprintf(resource.KAlbumTipTextTemplate, v.current.Count(), v.current.Date().Format(time.DateTime))

	//update music cards
	// keys := make([]model.AlbumKey, 0, len(album.MusicKeys))

	//filter title
	// for _, key := range keys {
	// 	if v.searchBar.Filter(v.client.Album(key).Title) {
	// 		keys = append(keys, key)
	// 	}
	// }

	//sort albums
	// slices.SortStableFunc(keys, func(a, b model.AlbumKey) int {
	// 	return v.searchBar.Sort(v.client.Album(a), v.client.Album(b))
	// })

	v.Refresh()
}
