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
	list *cwidget.SearchList[model.Music, *MusicCard]

	client   *model.Client
	pipeline DataPipeline[model.Music]
	current  model.Album
}

func NewMusicView(client *model.Client) *MusicView {
	var v *MusicView
	v = &MusicView{
		list: cwidget.NewSearchList(
			container.NewVBox(),
			newMusicCard,
			func(e cwidget.ItemTapEvent[model.Music]) {
				fmt.Println("play", e.Data.Title())
			},
			func(e cwidget.ItemTapEvent[model.Music]) {
				deleteMenu := cwidget.NewMenuItemWithIcon(resource.KDeleteText, theme.DeleteIcon(), func() { v.showDeleteMusicDialog(e.Data) })
				widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", deleteMenu), getWindow().Canvas(), e.AbsolutePosition)
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
		pipeline: DataPipeline[model.Music]{
			comparator: func(_, _ model.Music) int { return -1 },
			filter:     func(_ string) bool { return true },
		},
	}

	v.list.AddDropDown(cwidget.NewMenuItemWithIcon(resource.KMostRecentText, theme.HistoryIcon(), func() {
		v.pipeline.comparator = func(l, r model.Music) int {
			return -l.Date().Compare(r.Date())
		}
		v.updateList()
	}))
	v.list.AddDropDown(cwidget.NewMenuItemWithIcon(resource.KAlphabeticalText, resource.AlphabeticalIcon, func() {
		v.pipeline.comparator = func(l, r model.Music) int {
			return strings.Compare(strings.ToLower(l.Title()), strings.ToLower(r.Title()))
		}
		v.updateList()
	}))
	v.list.AddToolbar(cwidget.NewToolbarButton(resource.KReturnText, theme.NavigateBackIcon(), client.FocusAlbumView))
	v.ExtendBaseWidget(v)

	client.OnAlbumSelected().Attach(v)                                          //update current album and list content when selecting album
	client.OnAlbumsChanged().AttachFunc(func([]model.Album) { v.updateList() }) //update list content when albums get updated
	return v
}

func (v *MusicView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.list)
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
		v.list.ClearSearchEntry() //this triggers onChanged(), which triggers updateList()
	}
}

func (v *MusicView) updateList() {
	if v.current.Key().IsEmpty() {
		return
	}

	var err error
	v.current, err = v.client.GetAlbum(v.current.Key())
	if err != nil {
		fyne.LogError("client GetAlbum fails", err)
	}

	v.list.Update(v.pipeline.pass(v.current.Music()))
}
