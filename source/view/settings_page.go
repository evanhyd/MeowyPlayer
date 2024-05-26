package view

import (
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SettingPage struct {
	widget.BaseWidget
}

func NewSettingPage(client *model.Client) *SettingPage {
	v := &SettingPage{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *SettingPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Setting"))
}
