package view

import (
	"playground/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AccountPage struct {
	widget.BaseWidget
}

func NewAccountPage(client *model.Client) *AccountPage {
	v := &AccountPage{}
	v.ExtendBaseWidget(v)
	return v
}

func (v *AccountPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Account"))
}
