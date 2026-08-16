package ui

import (
	"log/slog"
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

type LoginPage struct {
	widget.BaseWidget
	userContext          *mcontext.UserContext
	userIdEntry          *widget.Entry
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	submitButton         *widget.Button
	goRegisterButton     *widget.Button
}

func newLoginPage(userContext *mcontext.UserContext, onGoRegister func()) *LoginPage {
	var p LoginPage
	p = LoginPage{
		userContext:   userContext,
		userIdEntry:   widget.NewEntry(),
		passwordEntry: widget.NewPasswordEntry(),
		submitButton: widget.NewButtonWithIcon(lang.L("Submit"), theme.LoginIcon(), func() {
			p.login(p.userIdEntry.Text, p.passwordEntry.Text)
		}),
		goRegisterButton: widget.NewButton(lang.L("Go Register"), onGoRegister),
	}
	p.ExtendBaseWidget(&p)

	checker := func(string) {
		if p.userIdEntry.Validate() == nil && p.passwordEntry.Validate() == nil {
			p.submitButton.Enable()
		} else {
			p.submitButton.Disable()
		}
	}

	p.userIdEntry.Validator = isValidUserId
	p.passwordEntry.Validator = isValidPassword
	p.userIdEntry.OnChanged = checker
	p.passwordEntry.OnChanged = checker

	p.submitButton.Disable()

	return &p
}

func (p *LoginPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mcontainer.NewCenter(0.65, 0.4, widget.NewForm(
		widget.NewFormItem("", widget.NewLabelWithStyle(lang.L("Login Page"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		widget.NewFormItem(lang.L("User ID"), p.userIdEntry),
		widget.NewFormItem(lang.L("Password"), p.passwordEntry),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.submitButton, p.goRegisterButton)),
	)))
}

func (p *LoginPage) clearEntry() {
	p.userIdEntry.SetText("")
	p.passwordEntry.SetText("")
}

func (p *LoginPage) login(userId string, password string) {
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
	})
}
