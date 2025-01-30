package view

import (
	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view/internal/cwidget"
	"meowyplayer/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

func RunApp() {
	//initialize locale
	resource.RegisterTranslation()

	//create main app
	mainApp := app.New()
	mainApp.SetIcon(resource.WindowIcon())
	mainApp.Settings().SetTheme(resource.NewVanillaTheme())

	//create main window
	window := mainApp.NewWindow(resource.WindowTitle())
	window.Resize(resource.KWindowSize)
	window.CenterOnScreen()
	window.SetContent(newMainPanel())

	//create system tray
	window.SetCloseIntercept(window.Hide)
	if desktop, ok := mainApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("",
			cwidget.NewMenuItem("", theme.MediaSkipPreviousIcon(), player.Instance().Prev),
			cwidget.NewMenuItem("", theme.RadioButtonCheckedIcon(), player.Instance().Play),
			cwidget.NewMenuItem("", theme.MediaSkipNextIcon(), player.Instance().Next),
			fyne.NewMenuItem("Show", window.Show),
		))
	}

	//run the client
	if err := model.NetworkClient().LoginWithConfig(); err != nil {
		fyne.LogError("failed to login with config", err)
	}
	window.ShowAndRun()
}

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}
