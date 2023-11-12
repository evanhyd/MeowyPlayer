package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/resource"
)

type AccountView struct {
	widget.BaseWidget
	accountName *widget.Label
	accountID   *widget.Label
	autoBackup  *widget.Check
}

func NewAccountView() *AccountView {
	profile := &AccountView{
		accountName: widget.NewLabel(""),
		accountID:   widget.NewLabel(""),
		autoBackup:  widget.NewCheck("Auto Backup", nil),
	}
	profile.ExtendBaseWidget(profile)
	return profile
}

func (v *AccountView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewVBox(v.accountName, v.accountID, v.autoBackup))
}

func (v *AccountView) SetAccount(account *resource.Account) {
	v.accountName.SetText(account.Name)
	v.accountID.SetText(account.ID)
}

func (v *AccountView) SetOnAutoBackup(callback func(bool)) {
	v.autoBackup.OnChanged = callback
}
