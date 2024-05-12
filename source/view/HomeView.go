package view

import (
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HomeView struct {
	widget.BaseWidget
}

func NewHomeView(client *model.MusicClient) *HomeView {
	v := &HomeView{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *HomeView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Home"))
}
