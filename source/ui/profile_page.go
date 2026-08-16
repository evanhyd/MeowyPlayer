package ui

import (
	"log/slog"
	"meowyplayer/mcontext"
	"meowyplayer/storages"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProfilePage struct {
	widget.BaseWidget
	userContext           *mcontext.UserContext
	userIdLabel           *widget.Label
	usernameLabel         *widget.Label
	registrationDateLabel *widget.Label
	logoutButton          *widget.Button
}

func newProfilePage(userContext *mcontext.UserContext) *ProfilePage {
	var p ProfilePage
	p = ProfilePage{
		userContext:           userContext,
		userIdLabel:           widget.NewLabel(""),
		usernameLabel:         widget.NewLabel(""),
		registrationDateLabel: widget.NewLabel(""),
		logoutButton:          widget.NewButtonWithIcon(lang.L("Logout"), theme.LogoutIcon(), p.logout),
	}

	userContext.AddListener(mcontext.OnPutUserEvent, func(data any) {
		p.setProfile(data.(mcontext.OnPutUserEventData).UserProfile)
	})

	userContext.AddListener(mcontext.OnSetStorageEvent, func(any) {
		profile, err := p.userContext.GetUser()
		if err != nil {
			slog.Error("failed to get user profile", "error", err)
		}
		p.setProfile(profile)
	})

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *ProfilePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(widget.NewForm(
		widget.NewFormItem(lang.L("User ID"), p.userIdLabel),
		widget.NewFormItem(lang.L("Username"), p.usernameLabel),
		widget.NewFormItem(lang.L("Registration Date"), p.registrationDateLabel),
		widget.NewFormItem("", container.NewHBox(layout.NewSpacer(), p.logoutButton)),
	))
}

func (p *ProfilePage) setProfile(profile storages.UserProfile) {
	p.userIdLabel.SetText(profile.UserId)
	p.usernameLabel.SetText(profile.Username)
	registrationDate := time.Unix(profile.RegistrationDate, 0)
	p.registrationDateLabel.SetText(registrationDate.Format(time.DateTime))
	p.userIdLabel.SetText(profile.UserId)
}

func (p *ProfilePage) logout() {
	if err := p.userContext.DeleteUser(); err != nil {
		slog.Error("failed to delete user", "error", err)
	}
}
