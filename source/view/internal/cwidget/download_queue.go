package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DownloadQueue struct {
	widget.BaseWidget
	button *widget.Button
	items  []fyne.Widget
}

func NewDownloadQueue(label string, icon fyne.Resource) *DownloadQueue {
	q := DownloadQueue{
		button: NewButton(label, icon, nil),
	}
	q.ExtendBaseWidget(&q)
	return &q
}

func (q *DownloadQueue) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(q.button)
}
