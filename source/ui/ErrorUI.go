package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func DisplayErrorIfNotNil(err error) {
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
	}
}
