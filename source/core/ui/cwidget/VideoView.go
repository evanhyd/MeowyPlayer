package cwidget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/utility/network/fileformat"
)

type VideoView struct {
	widget.BaseWidget
	tappableBase
	thumbnail *canvas.Image
	title     *widget.Label
	detail    *widget.Label
	download  *widget.Button
	highlight *canvas.Rectangle
}

func NewVideoView(result *fileformat.VideoResult, size fyne.Size, onDownload func()) *VideoView {
	const kConversionFactor = 60
	mins := int(result.Length.Minutes()) % kConversionFactor
	secs := int(result.Length.Seconds()) % kConversionFactor

	view := &VideoView{
		thumbnail: canvas.NewImageFromResource(result.Thumbnail),
		title:     widget.NewLabelWithStyle(fmt.Sprintf("[%02v:%02v] %v", mins, secs, result.Title), fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Symbol: true}),
		detail:    widget.NewLabel(result.ChannelTitle + "\n" + result.Stats),
		download:  NewButtonWithIcon("", theme.DownloadIcon(), nil),
		highlight: canvas.NewRectangle(theme.HoverColor()),
	}
	view.thumbnail.SetMinSize(size)
	view.title.Wrapping = fyne.TextWrapWord
	view.highlight.Hide()
	view.download.OnTapped = func() {
		view.download.Disable()
		view.download.SetIcon(theme.MoreHorizontalIcon())
		go func() {
			onDownload()
			view.download.SetIcon(theme.ConfirmIcon())
		}()
	}
	view.ExtendBaseWidget(view)
	return view
}

func (v *VideoView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		v.highlight,
		container.NewBorder(nil, nil, v.thumbnail, v.download, container.NewBorder(v.title, v.detail, nil, nil)),
	))
}

func (v *VideoView) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.Refresh()
}

func (v *VideoView) MouseOut() {
	v.highlight.Hide()
	v.Refresh()
}

func (v *VideoView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
