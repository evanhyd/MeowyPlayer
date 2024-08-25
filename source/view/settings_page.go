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

type LocalState struct {
	widget.BaseWidget
	hintLabel      *widget.Label
	loginButton    *widget.Button
	registerButton *widget.Button
}

func newLocalState() *LocalState {
	s := LocalState{
		hintLabel:      widget.NewLabel(resource.LoginToContinueText()),
		loginButton:    cwidget.NewButton(resource.LoginText(), nil, showLoginDialog),
		registerButton: cwidget.NewButton(resource.RegisterText(), nil, showRegisterDialog),
	}
	s.ExtendBaseWidget(&s)
	return &s
}

func (s *LocalState) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewVBox(s.hintLabel, s.loginButton, s.registerButton))
}

type RemoteState struct {
	widget.BaseWidget
	usernameLabel *widget.Label
	logoutButton  *widget.Button
}

func newRemoteState() *RemoteState {
	s := RemoteState{
		usernameLabel: widget.NewLabel(""),
		logoutButton:  cwidget.NewButton(resource.LogoutText(), nil, model.NetworkClient().Logout),
	}
	s.ExtendBaseWidget(&s)
	return &s
}

func (s *RemoteState) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewVBox(s.usernameLabel, s.logoutButton))
}

type SettingPage struct {
	widget.BaseWidget
	userstate fyne.CanvasObject
	vbox      *fyne.Container
}

func newSettingPage() *SettingPage {
	p := SettingPage{vbox: container.NewVBox()}
	p.setState(newLocalState())
	p.ExtendBaseWidget(&p)

	model.NetworkClient().OnConnectionChanged().Attach(&p)
	return &p
}

func (p *SettingPage) setState(state fyne.CanvasObject) {
	p.userstate = state
	p.vbox.RemoveAll()
	p.vbox.Add(state)
	p.Refresh()
}

func (p *SettingPage) Notify(isConnected bool) {
	if isConnected {
		p.setState(newRemoteState())
	} else {
		p.setState(newLocalState())
	}
}

func (p *SettingPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.vbox)
}
