package view

import (
	"meowyplayer/model"
	"meowyplayer/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func RunApp() {
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
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
	}

	//run the client
	if err := model.UIClient().Run(); err != nil {
		fyne.LogError("failed to run the UI client", err)
		return
	}
	if err := model.NetworkClient().Run(); err != nil {
		fyne.LogError("failed to run the network client", err)
		return
	}
	window.ShowAndRun()
}

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}
