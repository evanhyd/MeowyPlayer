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
			if err := model.NetworkClient().LoginManually(username.Text, password.Text); err != nil {
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

func showUploadLocalToTheAccountDialog() {
	dialog.ShowConfirm(resource.UploadLocalAlbumsToTheAccountText(), resource.MigrateConfirmationText(), func(yes bool) {
		if yes {
			progressBar := dialog.NewCustomWithoutButtons(resource.UploadText(), widget.NewProgressBarInfinite(), getWindow())
			progressBar.Show()
			defer progressBar.Hide()

			if err := model.NetworkClient().UploadLocalToTheAccount(); err != nil {
				fyne.LogError("failed to upload local albums to remote", err)
			}
		}
	}, getWindow())
}

func showBackupAlbumsToLocalDialog() {
	dialog.ShowConfirm(resource.BackupAlbumsToLocalText(), resource.MigrateConfirmationText(), func(yes bool) {
		if yes {
			progressBar := dialog.NewCustomWithoutButtons(resource.UploadText(), widget.NewProgressBarInfinite(), getWindow())
			progressBar.Show()
			defer progressBar.Hide()

			if err := model.NetworkClient().BackupAlbumsToLocal(); err != nil {
				fyne.LogError("failed to backup albums to local", err)
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
	usernameLabel         *widget.Label
	logoutButton          *widget.Button
	uploadToAccountButton *widget.Button
	backupToLocalButton   *widget.Button
}

func newRemoteState() *RemoteState {
	var s RemoteState
	s = RemoteState{
		usernameLabel:         widget.NewLabel(""),
		logoutButton:          cwidget.NewButton(resource.LogoutText(), nil, s.logout),
		uploadToAccountButton: cwidget.NewButton(resource.UploadLocalAlbumsToTheAccountText(), nil, showUploadLocalToTheAccountDialog),
		backupToLocalButton:   cwidget.NewButton(resource.BackupAlbumsToLocalText(), nil, showBackupAlbumsToLocalDialog),
	}
	s.ExtendBaseWidget(&s)
	return &s
}

func (s *RemoteState) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewVBox(s.usernameLabel, s.logoutButton, s.uploadToAccountButton, s.backupToLocalButton))
}

func (s *RemoteState) logout() {
	if err := model.NetworkClient().Logout(); err != nil {
		fyne.LogError("failed to logout", err)
	}
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

	model.NetworkClient().OnConnected().AttachFunc(func(info model.UserProfile) {
		p.setState(newRemoteState())
		p.userstate.(*RemoteState).usernameLabel.SetText(info.Username)
	})

	model.NetworkClient().OnDisconnected().AttachFunc(func(_ bool) {
		p.setState(newLocalState())
	})
	return &p
}

func (p *SettingPage) setState(state fyne.CanvasObject) {
	p.userstate = state
	p.vbox.RemoveAll()
	p.vbox.Add(state)
	p.Refresh()
}

func (p *SettingPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.vbox)
}
