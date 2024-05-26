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
)

func main() {
	//curl http://localhost/debug/pprof/heap -O profile.log
	//go tool pprof profile.log
	//go http.ListenAndServe("localhost:80", nil)

	//create model
	config := model.NewLocalFileSystem()
	client := model.NewClient(&config)

	//create main app
	mainApp := app.NewWithID(resource.KWindowTitle)
	fyne.SetCurrentApp(mainApp)
	mainApp.SetIcon(resource.WindowIcon)
	// application.Settings().SetTheme()

	//create main window
	window := mainApp.NewWindow(resource.KWindowTitle)
	window.SetCloseIntercept(window.Hide)
	window.CenterOnScreen()
	window.Resize(resource.KWindowSize)
	window.SetContent(view.NewMainPanel(&client))

	//create system tray
	if desktop, ok := mainApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
	}

	if err := client.Initialize(); err != nil {
		log.Println(err)
		return
	}

	window.ShowAndRun()
}
