package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/resource"
)

type MusicView struct {
	widget.BaseWidget
	tappableBase
	title     *widget.Label
	highlight *canvas.Rectangle
}

func NewMusicView(music *resource.Music) *MusicView {
	view := &MusicView{
		title:     widget.NewLabel(music.Description()),
		highlight: canvas.NewRectangle(theme.HoverColor()),
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
	//satisfy MouseMovement interface
}
