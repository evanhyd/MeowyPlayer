package cwidget

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/resource"
)

type CollectionInfoView struct {
	widget.BaseWidget
	tappableBase
	title     *widget.Label
	download  *widget.Button
	highlight *canvas.Rectangle
}

func NewCollectionInfoView(info *resource.CollectionInfo, onDownload func(*resource.CollectionInfo)) *CollectionInfoView {
	view := &CollectionInfoView{
		title:     widget.NewLabel(fmt.Sprintf("%v %v %.1f mb", info.Title, info.Date.Format(time.DateTime), float64(info.Size)/1024/1024)),
		download:  NewButtonWithIcon("", theme.DownloadIcon(), nil),
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	view.highlight.Hide()
	view.download.OnTapped = func() {
		view.download.Disable()
		view.download.SetIcon(theme.MoreHorizontalIcon())
		onDownload(info)
		view.download.SetIcon(theme.DownloadIcon())
		view.download.Enable()
	}
	view.ExtendBaseWidget(view)
	return view
}

func (v *CollectionInfoView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		v.highlight,
		container.NewBorder(
			nil,
			nil,
			nil,
			v.download,
			v.title,
		),
	))
}

func (v *CollectionInfoView) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.Refresh()
}

func (v *CollectionInfoView) MouseOut() {
	v.highlight.Hide()
	v.Refresh()
}

func (v *CollectionInfoView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
