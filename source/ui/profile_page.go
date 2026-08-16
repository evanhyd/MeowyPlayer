package ui

import (
	"errors"
	"log/slog"
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
		return errors.New("user id length must be between 3 - 20 characters long")
	}

	hasLetter := false
	for _, char := range userId {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("user id must contain only letter and optionally number")
		}
		if unicode.IsLetter(char) {
			hasLetter = true
		}
	}

	if !hasLetter {
		return errors.New("user id must contain only letter and optionally number")
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
	userIdEntry          *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	loginButton          *widget.Button
	registerButton       *widget.Button
	onLogin              func()
}

func newLoginPage(userContext *mcontext.UserContext, onLogin func(), onRegister func()) *loginPage {
	var p loginPage
	p = loginPage{
		userContext:   userContext,
		userIdEntry:   widget.NewEntry(),
		passwordEntry: widget.NewPasswordEntry(),
		loginButton: widget.NewButton(lang.L("Login"), func() {
			p.login(p.userIdEntry.Text, p.passwordEntry.Text)
		}),
		registerButton: widget.NewButton(lang.L("Go Register"), onRegister),
		onLogin:        onLogin,
	}
	p.ExtendBaseWidget(&p)

	checker := func(string) {
		if p.userIdEntry.Validate() == nil && p.passwordEntry.Validate() == nil {
			p.loginButton.Enable()
		} else {
			p.loginButton.Disable()
		}
	}

	p.userIdEntry.Validator = isValidUserId
	p.passwordEntry.Validator = isValidPassword
	p.userIdEntry.OnChanged = checker
	p.passwordEntry.OnChanged = checker

	p.loginButton.Disable()

	return &p
}

func (p *loginPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mcontainer.NewCenter(0.65, 0.3, widget.NewForm(
		widget.NewFormItem(lang.L("UserId"), p.userIdEntry),
		widget.NewFormItem(lang.L("Password"), p.passwordEntry),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.loginButton, p.registerButton)),
	)))
}

func (p *loginPage) login(userId string, password string) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	waitingDialog := dialog.NewCustomWithoutButtons(lang.L("logining"), widget.NewProgressBarInfinite(), win)

	fyne.Do(func() {
		waitingDialog.Show()
		err := func() error {
			// Login to get token.
			token, err := func() (string, error) {
				request := handlers.LoginRequest{
					UserId:   userId,
					Password: password,
				}
				response := handlers.LoginResponse{}
				err := handlers.SendJSON(p.userContext.Config().Endpoints["login"], request, &response)
				return response.Token, err
			}()
			if err != nil {
				slog.Error("failed to login", "error", err)
				return err
			}

			// Get user profile.
			profile, err := func() (storages.UserProfile, error) {
				request := handlers.MeRequest{
					Token: token,
				}
				response := handlers.MeResponse{}
				err := handlers.SendJSON(p.userContext.Config().Endpoints["me"], request, &response)

				profile := storages.UserProfile{
					UserId:           response.UserId,
					Username:         response.Username,
					Language:         response.Language,
					RegistrationDate: response.RegistrationDate,
					Token:            token,
				}
				return profile, err
			}()
			if err != nil {
				slog.Error("failed to get user profile", "error", err)
				return err
			}

			if err := p.userContext.PutUser(profile); err != nil {
				slog.Error("failed to put user profile in the storage", "error", err)
				return err
			}
			return nil
		}()

		waitingDialog.Dismiss()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		p.onLogin()
	})
}

type registerPage struct {
	widget.BaseWidget
	userContext          *mcontext.UserContext
	userIdEntry          *widget.Entry
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
		userIdEntry:          widget.NewEntry(),
		passwordEntry:        widget.NewPasswordEntry(),
		confirmPasswordEntry: widget.NewPasswordEntry(),
		submitButton: widget.NewButton(lang.L("Submit"), func() {
			p.registerUser(p.userIdEntry.Text, p.passwordEntry.Text)
		}),
		cancelButton: widget.NewButton(lang.L("Cancel"), onCancel),
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
			return errors.New("password is not the same")
		}
		return nil
	}
	p.userIdEntry.OnChanged = checker
	p.passwordEntry.OnChanged = checker
	p.confirmPasswordEntry.OnChanged = checker
	p.submitButton.Disable()

	return &p
}

func (p *registerPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mcontainer.NewCenter(0.65, 0.3, widget.NewForm(
		widget.NewFormItem(lang.L("User ID"), p.userIdEntry),
		widget.NewFormItem(lang.L("Password"), p.passwordEntry),
		widget.NewFormItem(lang.L("Confirm Password"), p.confirmPasswordEntry),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.submitButton, p.cancelButton)),
	)))
}

func (p *registerPage) registerUser(userId string, password string) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	waitingDialog := dialog.NewCustomWithoutButtons(lang.L("registering"), widget.NewProgressBarInfinite(), win)

	fyne.Do(func() {
		waitingDialog.Show()

		// Send request..
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

		// Remove the old entry values.
		p.userIdEntry.SetText("")
		p.passwordEntry.SetText("")
		p.confirmPasswordEntry.SetText("")

		// Go back to the login page.
		dialog.ShowInformation(lang.L("register"), lang.L("Register successfully!"), win)
		p.cancelButton.OnTapped()
	})
}

type UserProfilePage struct {
	widget.BaseWidget
	userContext           *mcontext.UserContext
	userIdLabel           *widget.Label
	usernameLabel         *widget.Label
	registrationDateLabel *widget.Label
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
		userContext: userContext,
		loginView: newLoginPage(
			userContext,
			func() {
				dialog.ShowInformation(lang.L("login"), lang.L("Welcome back"), fyne.CurrentApp().Driver().AllWindows()[0])
			},
			p.showRegisterPage,
		),
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
