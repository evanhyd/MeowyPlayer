package view

import (
	"log"
	"playground/model"
	"playground/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

func RunApp() {
	//create main app
	mainApp := app.NewWithID(resource.KWindowTitle)
	fyne.SetCurrentApp(mainApp)
	mainApp.SetIcon(resource.WindowIcon)
	mainApp.Settings().SetTheme(theme.DarkTheme())

	//create main window
	window := mainApp.NewWindow(resource.KWindowTitle)
	window.SetCloseIntercept(window.Hide)
	window.CenterOnScreen()
	window.Resize(resource.KWindowSize)
	window.SetContent(newMainPanel())

	//create system tray
	if desktop, ok := mainApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
	}

	//run the client
	if err := model.Instance().Run(); err != nil {
		log.Println(err)
		return
	}
	window.ShowAndRun()
}

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}
