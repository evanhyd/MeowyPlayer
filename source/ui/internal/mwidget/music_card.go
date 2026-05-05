package mwidget

import (
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = &MusicCard{}

type MusicCard struct {
	widget.BaseWidget
	highlight   *canvas.Rectangle
	title       *widget.Label
	description *widget.Label
	music       storages.Music
	onTapped    func(playlist storages.Music)
}

func NewMusicCard(onTapped func(music storages.Music)) *MusicCard {
	c := MusicCard{
		highlight:   canvas.NewRectangle(theme.Color(theme.ColorNameHover)),
		title:       widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		description: widget.NewLabel(""),
		onTapped:    onTapped,
	}
	c.highlight.Hide()
	c.highlight.CornerRadius = 8.0
	c.title.Truncation = fyne.TextTruncateEllipsis
	c.title.Wrapping = fyne.TextWrapBreak
	c.ExtendBaseWidget(&c)
	return &c
}

func (c *MusicCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		c.highlight,
		container.NewBorder(nil, nil, nil, c.description, c.title),
	))
}

func (c *MusicCard) MouseIn(*desktop.MouseEvent) {
	c.highlight.Show()
	c.Refresh()
}

func (c *MusicCard) MouseOut() {
	c.highlight.Hide()
	c.Refresh()
}

func (c *MusicCard) MouseMoved(*desktop.MouseEvent) {
	// Satisfy desktop.Hoverable
}

func (c *MusicCard) Tapped(*fyne.PointEvent) {
	c.onTapped(c.music)
}

func (c *MusicCard) Set(music storages.Music) {
	c.title.SetText(music.Title)
	c.description.SetText(mutil.SecondsToTime(music.LengthSeconds))
	c.music = music
	c.Refresh()
}
