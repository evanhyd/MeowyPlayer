package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type MusicView struct {
	*widget.Label
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func NewMusicView(music *player.Music) *MusicView {
	view := &MusicView{Label: widget.NewLabel(music.Description()), OnTapped: func(*fyne.PointEvent) {}, OnTappedSecondary: func(*fyne.PointEvent) {}}
	view.ExtendBaseWidget(view)
	return view
}

func (m *MusicView) Tapped(event *fyne.PointEvent) {
	m.OnTapped(event)
}

func (m *MusicView) TappedSecondary(event *fyne.PointEvent) {
	m.OnTappedSecondary(event)
}
