package ui

import (
	"errors"
	"meowyplayer/ui/internal/mcontainer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type loginPage struct {
	widget.BaseWidget
	usernameEntry        *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	loginButton          *widget.Button
	registerButton       *widget.Button
}

func newLoginPage(onRegister func()) *loginPage {
	var p loginPage
	p = loginPage{
		usernameEntry:  widget.NewEntry(),
		passwordEntry:  widget.NewPasswordEntry(),
		loginButton:    widget.NewButton(lang.L("Login"), nil),
		registerButton: widget.NewButton(lang.L("Register"), onRegister),
	}
	p.ExtendBaseWidget(&p)
	return &p
}

func (p *loginPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mcontainer.NewCenter(0.65, 0.3, widget.NewForm(
		widget.NewFormItem(lang.L("Username"), p.usernameEntry),
		widget.NewFormItem(lang.L("Password"), p.passwordEntry),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.loginButton, p.registerButton)),
	)))
}

type registerPage struct {
	widget.BaseWidget
	usernameEntry        *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	captchaEntry         *widget.Entry
	submitButton         *widget.Button
	cancelButton         *widget.Button
}

func newRegisterPage(onCancel func()) *registerPage {
	p := registerPage{
		usernameEntry:        widget.NewEntry(),
		passwordEntry:        widget.NewPasswordEntry(),
		confirmPasswordEntry: widget.NewPasswordEntry(),
		submitButton:         widget.NewButton(lang.L("Submit"), nil),
		cancelButton:         widget.NewButton(lang.L("Cancel"), onCancel),
	}
	p.ExtendBaseWidget(&p)

	p.confirmPasswordEntry.AlwaysShowValidationError = true
	p.confirmPasswordEntry.Validator = func(s string) error {
		if s != p.passwordEntry.Text {
			return errors.New("not matching with the password")
		}
		return nil
	}

	return &p
}

func (p *registerPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mcontainer.NewCenter(0.65, 0.3, widget.NewForm(
		widget.NewFormItem(lang.L("Username"), p.usernameEntry),
		widget.NewFormItem(lang.L("Password"), p.passwordEntry),
		widget.NewFormItem(lang.L("Confirm Password"), p.confirmPasswordEntry),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.submitButton, p.cancelButton)),
	)))
}

type ProfilePage struct {
	widget.BaseWidget
	loginView    *loginPage
	registerView *registerPage
}

func newProfilePage() *ProfilePage {
	var p ProfilePage
	p = ProfilePage{
		loginView:    newLoginPage(p.showRegisterPage),
		registerView: newRegisterPage(p.showLoginPage),
	}
	p.ExtendBaseWidget(&p)

	p.registerView.Hide()
	return &p
}

func (p *ProfilePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(p.registerView, p.loginView))
}

func (p *ProfilePage) showLoginPage() {
	p.registerView.Hide()
	p.loginView.Show()
}

func (p *ProfilePage) showRegisterPage() {
	p.loginView.Hide()
	p.registerView.Show()
}
