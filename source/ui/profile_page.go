package ui

import (
	"errors"
	"meowyplayer/handlers"
	"meowyplayer/mcontext"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"unicode"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func isValidUserId(userId string) error {
	length := utf8.RuneCountInString(userId)
	if length < 3 || length > 20 {
		return errors.New("username length must be between 3 - 20 characters long")
	}

	hasLetter := false
	for _, char := range userId {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("username must contain only letter and optionally number")
		}
		if unicode.IsLetter(char) {
			hasLetter = true
		}
	}

	if !hasLetter {
		return errors.New("username must contain only letter and optionally number")
	}
	return nil
}

// 8-30 characters long and contains at least one letter, one number, and one special character.
func isValidPassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < 8 || length > 30 {
		return errors.New("password length must be between 8 - 30 characters long")
	}

	var hasLetter, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char): // Matches digits in many scripts
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !(hasLetter && hasNumber && hasSpecial) {
		return errors.New("password must contain at least letter, number, and a special character !@#$%^&*()")
	}
	return nil
}

type loginPage struct {
	widget.BaseWidget
	userContext          *mcontext.UserContext
	usernameEntry        *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	loginButton          *widget.Button
	registerButton       *widget.Button
}

func newLoginPage(userContext *mcontext.UserContext, onRegister func()) *loginPage {
	var p loginPage
	p = loginPage{
		userContext:    userContext,
		usernameEntry:  widget.NewEntry(),
		passwordEntry:  widget.NewPasswordEntry(),
		loginButton:    widget.NewButton(lang.L("Login"), nil),
		registerButton: widget.NewButton(lang.L("Go Register"), onRegister),
	}
	p.ExtendBaseWidget(&p)

	checker := func(string) {
		if p.usernameEntry.Validate() == nil && p.passwordEntry.Validate() == nil {
			p.loginButton.Enable()
		} else {
			p.loginButton.Disable()
		}
	}

	p.usernameEntry.Validator = isValidUserId
	p.passwordEntry.Validator = isValidPassword
	p.usernameEntry.OnChanged = checker
	p.passwordEntry.OnChanged = checker

	p.loginButton.Disable()

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
	userContext          *mcontext.UserContext
	usernameEntry        *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	captchaEntry         *widget.Entry
	submitButton         *widget.Button
	cancelButton         *widget.Button
}

func newRegisterPage(userContext *mcontext.UserContext, onCancel func()) *registerPage {
	var p registerPage
	p = registerPage{
		userContext:          userContext,
		usernameEntry:        widget.NewEntry(),
		passwordEntry:        widget.NewPasswordEntry(),
		confirmPasswordEntry: widget.NewPasswordEntry(),
		submitButton: widget.NewButton(lang.L("Submit"), func() {
			win := fyne.CurrentApp().Driver().AllWindows()[0]
			_, err := p.registerUser(p.usernameEntry.Text, p.passwordEntry.Text)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			dialog.ShowInformation(lang.L("register"), lang.L("Register successfully!"), win)
		}),
		cancelButton: widget.NewButton(lang.L("Cancel"), onCancel),
	}
	p.ExtendBaseWidget(&p)

	checker := func(string) {
		if p.usernameEntry.Validate() == nil && p.passwordEntry.Validate() == nil && p.confirmPasswordEntry.Validate() == nil {
			p.submitButton.Enable()
		} else {
			p.submitButton.Disable()
		}
	}

	p.usernameEntry.Validator = isValidUserId
	p.passwordEntry.Validator = isValidPassword
	p.confirmPasswordEntry.Validator = func(s string) error {
		if s != p.passwordEntry.Text {
			return errors.New("password is not the same")
		}
		return nil
	}
	p.usernameEntry.OnChanged = checker
	p.passwordEntry.OnChanged = checker
	p.confirmPasswordEntry.OnChanged = checker
	p.submitButton.Disable()

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

func (p *registerPage) registerUser(username string, password string) (handlers.RegisterResponse, error) {
	// Send request to the server.
	request := handlers.RegisterRequest{
		UserId:   username,
		Username: username,
		Password: password,
		Language: storages.LangEnglish,
	}

	response := handlers.RegisterResponse{}
	err := handlers.SendJSON(p.userContext.Config().Endpoints["register"], request, &response)
	return response, err

}

type ProfilePage struct {
	widget.BaseWidget
	userContext  *mcontext.UserContext
	loginView    *loginPage
	registerView *registerPage
}

func newProfilePage(userContext *mcontext.UserContext) *ProfilePage {
	var p ProfilePage
	p = ProfilePage{
		userContext:  userContext,
		loginView:    newLoginPage(userContext, p.showRegisterPage),
		registerView: newRegisterPage(userContext, p.showLoginPage),
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
