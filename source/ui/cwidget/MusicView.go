package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type MusicView struct {
	widget.BaseWidget
	title             *widget.Label
	highlight         *canvas.Rectangle
	OnTapped          func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func NewMusicView(music *player.Music) *MusicView {
	view := &MusicView{
		widget.BaseWidget{},
		widget.NewLabel(music.Description()),
		canvas.NewRectangle(theme.HoverColor()),
		func(*fyne.PointEvent) {},
		func(*fyne.PointEvent) {},
	}
	view.highlight.Hide()
	view.ExtendBaseWidget(view)
	return view
}

func (m *MusicView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewMax(m.title, m.highlight))
}

func (m *MusicView) MouseIn(event *desktop.MouseEvent) {
	m.highlight.Show()
	m.Refresh()
}

func (m *MusicView) MouseOut() {
	m.highlight.Hide()
	m.Refresh()
}

func (m *MusicView) MouseMoved(*desktop.MouseEvent) {
}

func (m *MusicView) Tapped(event *fyne.PointEvent) {
	m.OnTapped(event)
}

func (m *MusicView) TappedSecondary(event *fyne.PointEvent) {
	m.OnTappedSecondary(event)
}
