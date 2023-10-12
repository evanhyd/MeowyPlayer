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

type VideoResultView struct {
	widget.BaseWidget
	tappableBase
	thumbnail    *canvas.Image
	title        *widget.Label
	channelTitle *widget.Label
	stats        *widget.Label
	highlight    *canvas.Rectangle
}

func NewVideoResultView(result *fileformat.VideoResult, size fyne.Size) *VideoResultView {
	const kConversionFactor = 60
	mins := int(result.Length.Minutes()) % kConversionFactor
	secs := int(result.Length.Seconds()) % kConversionFactor

	view := &VideoResultView{
		thumbnail:    canvas.NewImageFromResource(result.Thumbnail),
		title:        widget.NewLabelWithStyle(fmt.Sprintf("[%02v:%02v] %v", mins, secs, result.Title), fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Symbol: true}),
		channelTitle: widget.NewLabel(result.ChannelTitle),
		stats:        widget.NewLabel(result.Stats),
		highlight:    canvas.NewRectangle(theme.HoverColor()),
	}
	view.thumbnail.SetMinSize(size)
	view.title.Wrapping = fyne.TextWrapWord
	view.highlight.Hide()
	view.ExtendBaseWidget(view)
	return view
}

func (v *VideoResultView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewMax(
		v.highlight,
		container.NewBorder(
			nil,
			nil,
			v.thumbnail,
			nil,
			container.NewVBox(v.title, v.channelTitle, v.stats),
		),
	))
}

func (v *VideoResultView) MouseIn(*desktop.MouseEvent) {
	v.highlight.Show()
	v.Refresh()
}

func (v *VideoResultView) MouseOut() {
	v.highlight.Hide()
	v.Refresh()
}

func (v *VideoResultView) MouseMoved(*desktop.MouseEvent) {
	//satisfy MouseMovement interface
}
