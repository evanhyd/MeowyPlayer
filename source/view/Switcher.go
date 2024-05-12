package view

import (
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Switcher struct {
	widget.BaseWidget
	albumView fyne.CanvasObject
	musicView fyne.CanvasObject
}

func NewSwitcher(client *model.MusicClient) *Switcher {
	s := &Switcher{albumView: NewAlbumView(client), musicView: NewMusicView(client)}
	s.musicView.Hide()
	client.OnAlbumsChanged().AttachFunc(s.showAlbumTab)
	client.OnAlbumFocused().AttachFunc(s.showMusicTab)
	s.ExtendBaseWidget(s)
	return s
}

func (s *Switcher) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(s.albumView, s.musicView))
}

func (s *Switcher) showAlbumTab(_ []model.Album) {
	s.albumView.Show()
	s.musicView.Hide()
}

func (s *Switcher) showMusicTab(_ model.Album) {
	s.albumView.Hide()
	s.musicView.Show()
}
