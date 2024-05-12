package view

import (
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SettingView struct {
	widget.BaseWidget
}

func NewSettingView(client *model.MusicClient) *SettingView {
	v := &SettingView{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *SettingView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Setting"))
}
