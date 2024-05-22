package view

import (
	"fmt"
	"playground/model"
	"playground/resource"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type AlbumCard struct {
	widget.BaseWidget
	cover     *canvas.Image
	title     *widget.Label
	tip       *widget.Label
	isHovered bool
}

func newAlbumCard() *AlbumCard {
	v := &AlbumCard{
		cover: canvas.NewImageFromResource(nil),
		title: widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
		tip:   widget.NewLabel(""),
	}

	v.cover.FillMode = canvas.ImageFillContain
	v.title.Truncation = fyne.TextTruncateEllipsis
	v.tip.Wrapping = fyne.TextWrapWord
	v.tip.Hide()

	v.ExtendBaseWidget(v)
	return v
}

func (v *AlbumCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		container.NewBorder(nil, v.title, nil, nil, v.cover),
		v.tip,
	))
}

func (v *AlbumCard) MouseIn(e *desktop.MouseEvent) {
	v.cover.Translucency = 0.8
	v.cover.Refresh()
	v.tip.Show()
	v.isHovered = true
}

func (v *AlbumCard) MouseOut() {
	v.cover.Translucency = 0.0
	v.cover.Refresh()
	v.tip.Hide()
	v.isHovered = false
}

func (v *AlbumCard) MouseMoved(*desktop.MouseEvent) {
	//satisfy Hoverable interface
}

func (v *AlbumCard) Cursor() desktop.Cursor {
	if v.isHovered {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

func (v *AlbumCard) Notify(album model.Album) {
	v.cover.Resource = album.Cover()
	v.cover.Refresh()
	v.title.SetText(album.Title())
	v.tip.SetText(fmt.Sprintf(resource.KAlbumTipTextTemplate, album.Count(), album.Date().Format(time.DateTime)))
}
