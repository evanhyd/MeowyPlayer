package ui

import (
	"fyne.io/fyne/v2/dialog"
	"meowyplayer.com/source/player"
)

func DisplayErrorIfAny(err error) {
	if err != nil {
		dialog.ShowError(err, player.GetMainWindow())
	}
}
