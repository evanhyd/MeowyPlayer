package cwidget

import (
	"fmt"
	"playground/model"
	"playground/view/internal/resource"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type AlbumCard struct {
	widget.BaseWidget
	TappableComponent
	CursorableComponent
	cover *canvas.Image
	title *widget.Label
	tip   *widget.Label
	key   model.AlbumKey
}

func NewAlbumCardConstructor(onTapped func(model.AlbumKey), onTappedSecondary func(*fyne.PointEvent, model.AlbumKey)) func() *AlbumCard {
	return func() *AlbumCard {
		v := AlbumCard{
			cover: canvas.NewImageFromResource(nil),
			title: widget.NewLabel(""),
			tip:   widget.NewLabel(""),
		}
		v.OnTapped = func(*fyne.PointEvent) { onTapped(v.key) }
		v.OnTappedSecondary = func(e *fyne.PointEvent) { onTappedSecondary(e, v.key) }
		v.ExtendBaseWidget(&v)
		return &v
	}
}

func (v *AlbumCard) CreateRenderer() fyne.WidgetRenderer {
	v.cover.FillMode = canvas.ImageFillContain
	v.title.Truncation = fyne.TextTruncateEllipsis
	v.title.Alignment = fyne.TextAlignCenter
	v.tip.Wrapping = fyne.TextWrapWord
	v.tip.Hide()

	return widget.NewSimpleRenderer(container.NewStack(container.NewBorder(nil, v.title, nil, nil, v.cover), v.tip))
}

func (v *AlbumCard) MouseIn(e *desktop.MouseEvent) {
	v.cover.Translucency = 0.8
	v.cover.Refresh()
	v.tip.Show()
}

func (v *AlbumCard) MouseOut() {
	v.cover.Translucency = 0.0
	v.cover.Refresh()
	v.tip.Hide()
}

func (v *AlbumCard) MouseMoved(*desktop.MouseEvent) {
	//Hoverable interface
}

func (v *AlbumCard) Notify(album model.Album) {
	v.key = album.Key()
	v.cover.Resource = album.Cover()
	v.cover.Refresh()
	v.title.SetText(album.Title())
	v.tip.SetText(fmt.Sprintf(resource.KAlbumTipTextTemplate, album.Count(), album.Date().Format(time.DateTime)))
}
