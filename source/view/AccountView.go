package view

import (
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AccountView struct {
	widget.BaseWidget
}

func NewAccountView(client *model.MusicClient) *AccountView {
	v := &AccountView{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *AccountView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Account"))
}
