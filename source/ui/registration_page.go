package ui

import (
	"errors"
	"meowyplayer/handlers"
	"meowyplayer/mcontext"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type RegistrationPage struct {
	widget.BaseWidget
	userContext          *mcontext.UserContext
	userIdEntry          *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	captchaEntry         *widget.Entry
	submitButton         *widget.Button
	goLoginButton        *widget.Button
}

func newRegistrationPage(userContext *mcontext.UserContext, onGoLogin func()) *RegistrationPage {
	var p RegistrationPage
	p = RegistrationPage{
		userContext:          userContext,
		userIdEntry:          widget.NewEntry(),
		passwordEntry:        widget.NewPasswordEntry(),
		confirmPasswordEntry: widget.NewPasswordEntry(),
		submitButton: widget.NewButtonWithIcon(lang.L("Submit"), theme.MailSendIcon(), func() {
			p.registerUser(p.userIdEntry.Text, p.passwordEntry.Text)
		}),
		goLoginButton: widget.NewButton(lang.L("Go Login"), onGoLogin),
	}
	p.ExtendBaseWidget(&p)

	checker := func(string) {
		if p.userIdEntry.Validate() == nil && p.passwordEntry.Validate() == nil && p.confirmPasswordEntry.Validate() == nil {
			p.submitButton.Enable()
		} else {
			p.submitButton.Disable()
		}
	}

	p.userIdEntry.Validator = isValidUserId
	p.passwordEntry.Validator = isValidPassword
	p.confirmPasswordEntry.Validator = func(s string) error {
		if s != p.passwordEntry.Text {
			return errors.New(lang.L("password is not the same"))
		}
		return nil
	}
	p.userIdEntry.OnChanged = checker
	p.passwordEntry.OnChanged = checker
	p.confirmPasswordEntry.OnChanged = checker
	p.submitButton.Disable()

	return &p
}

func (p *RegistrationPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mcontainer.NewCenter(0.65, 1.0, widget.NewForm(
		widget.NewFormItem("", widget.NewLabelWithStyle(lang.L("Registration Page"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		widget.NewFormItem(lang.L("User ID"), p.userIdEntry),
		widget.NewFormItem(lang.L("Password"), p.passwordEntry),
		widget.NewFormItem(lang.L("Confirm Password"), p.confirmPasswordEntry),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.submitButton, p.goLoginButton)),
	)))
}

func (p *RegistrationPage) clearEntry() {
	p.userIdEntry.SetText("")
	p.passwordEntry.SetText("")
	p.confirmPasswordEntry.SetText("")
}

func (p *RegistrationPage) registerUser(userId string, password string) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	waitingDialog := dialog.NewCustomWithoutButtons(lang.L("registering"), widget.NewProgressBarInfinite(), win)

	fyne.Do(func() {
		waitingDialog.Show()

		// Send registration request.
		request := handlers.RegisterRequest{
			UserId:   userId,
			Username: userId,
			Password: password,
			Language: storages.LangEnglish,
		}

		response := handlers.RegisterResponse{}
		err := handlers.SendJSON(p.userContext.Config().Endpoints["register"], request, &response)
		waitingDialog.Dismiss()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		// Go back to the login page.
		dialog.ShowInformation(lang.L("register"), lang.L("Register successfully!"), win)
		p.goLoginButton.OnTapped()
	})
}
