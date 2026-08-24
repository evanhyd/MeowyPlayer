package ui

import (
	"errors"
	"meowyplayer/mcontext"
	"unicode"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type AuthPage struct {
	widget.BaseWidget
	loginPage        *LoginPage
	registrationPage *RegistrationPage
}

func newAuthPage(userContext *mcontext.UserContext) *AuthPage {
	var p AuthPage
	p = AuthPage{
		loginPage:        newLoginPage(userContext, p.showRegisterPage),
		registrationPage: newRegistrationPage(userContext, p.showLoginPage),
	}
	p.ExtendBaseWidget(&p)
	p.showLoginPage()
	return &p
}

func (p *AuthPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(p.registrationPage, p.loginPage))
}

func (p *AuthPage) showLoginPage() {
	p.loginPage.clearEntry()
	p.registrationPage.Hide()
	p.loginPage.Show()
}

func (p *AuthPage) showRegisterPage() {
	p.registrationPage.clearEntry()
	p.loginPage.Hide()
	p.registrationPage.Show()
}

func isValidUserId(userId string) error {
	length := utf8.RuneCountInString(userId)
	if length < 3 || length > 20 {
		return errors.New(lang.L("user id length must be between 3 - 20 characters long"))
	}

	hasLetter := false
	for _, char := range userId {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New(lang.L("user id must contain only letter and optionally number"))
		}
		if unicode.IsLetter(char) {
			hasLetter = true
		}
	}

	if !hasLetter {
		return errors.New(lang.L("user id must contain only letter and optionally number"))
	}
	return nil
}

// 8-30 characters long and contains at least one letter, one number, and one special character.
func isValidPassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < 8 || length > 30 {
		return errors.New(lang.L("password length must be between 8 - 30 characters long"))
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
		return errors.New(lang.L("password must contain at least letter, number, and a special character !@#$%^&*()"))
	}
	return nil
}
