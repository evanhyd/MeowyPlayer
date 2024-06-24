package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SettingPage struct {
	widget.BaseWidget
}

func newSettingPage() *SettingPage {
	v := SettingPage{}
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *SettingPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Setting"))
}
