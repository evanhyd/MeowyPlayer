package main

import (
	"log"
	_ "net/http/pprof"

	"playground/model"
	"playground/resource"
	"playground/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

func main() {
	//curl http://localhost/debug/pprof/heap -O profile.log
	//go tool pprof profile.log
	//go http.ListenAndServe("localhost:80", nil)

	model.CreateClient(model.NewLocalStorage())

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
	window.SetContent(view.NewMainPanel())

	//create system tray
	if desktop, ok := mainApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
	}

	//run the client
	if err := model.GetClient().Run(); err != nil {
		log.Println(err)
		return
	}
	window.ShowAndRun()
}
