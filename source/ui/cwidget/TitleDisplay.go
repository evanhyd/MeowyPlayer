package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type titleDisplay struct {
	widget.BaseWidget
	title *widget.Label
}

func newTitleDisplay() *titleDisplay {
	display := &titleDisplay{widget.BaseWidget{}, widget.NewLabel("")}
	display.ExtendBaseWidget(display)
	return display
}

func (t *titleDisplay) CreateRenderer() fyne.WidgetRenderer {
	return t.title.CreateRenderer()
}

func (t *titleDisplay) SetMusicTitle(music *player.Music) {
	t.title.SetText(music.Title[:len(music.Title)-3])
}
