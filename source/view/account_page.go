package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AccountPage struct {
	widget.BaseWidget
}

func newAccountPage() *AccountPage {
	v := AccountPage{}
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *AccountPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewLabel("Account"))
}
