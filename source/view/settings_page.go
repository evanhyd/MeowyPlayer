package view

import (
	"meowyplayer/model"
	"meowyplayer/view/internal/cwidget"
	"meowyplayer/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showLoginDialog() {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem(resource.UsernameText(), username),
		widget.NewFormItem(resource.PasswordText(), password),
	}

	dialog.ShowForm(resource.LoginText(), resource.LoginText(), resource.CancelText(), items, func(login bool) {
		if login {
			if err := model.NetworkClient().Login(username.Text, password.Text); err != nil {
				dialog.ShowError(err, getWindow())
			}
		}
	}, getWindow())
}

func showRegisterDialog() {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem(resource.UsernameText(), username),
		widget.NewFormItem(resource.PasswordText(), password),
	}

	dialog.ShowForm(resource.RegisterText(), resource.RegisterText(), resource.CancelText(), items, func(login bool) {
		if login {
			if err := model.NetworkClient().Register(username.Text, password.Text); err != nil {
				dialog.ShowError(err, getWindow())
			}
		}
	}, getWindow())
}

type SettingPage struct {
	widget.BaseWidget
	index          *widget.RichText
	loginButton    *widget.Button
	registerButton *widget.Button
}

func newSettingPage() *SettingPage {
	v := SettingPage{
		index:          widget.NewRichTextWithText("Login to continue!"),
		loginButton:    cwidget.NewButton(resource.LoginText(), nil, showLoginDialog),
		registerButton: cwidget.NewButton(resource.RegisterText(), nil, showRegisterDialog),
	}
	v.ExtendBaseWidget(&v)
	return &v
}

func (v *SettingPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewVBox(v.index, v.loginButton, v.registerButton))
}
