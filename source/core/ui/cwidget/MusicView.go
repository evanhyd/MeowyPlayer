package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/resource"
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

func (v *MusicView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(v.title, v.highlight))
}

func (v *MusicView) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.Refresh()
}

func (v *MusicView) MouseOut() {
	v.highlight.Hide()
	v.Refresh()
}

func (v *MusicView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
