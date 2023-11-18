package main

import (
	"fmt"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui"
	"meowyplayer.com/utility/logger"
)

func main() {
	// redirect panic message
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Errorf("%v\n%v", err, string(debug.Stack())), 1)
		}
	}()

	logger.Initiate()
	resource.MakeNecessaryPath()

	//initiate app configuration
	fyne.SetCurrentApp(app.NewWithID("MeowyPlayer"))
	fyne.CurrentApp().Settings().SetTheme(resource.VanillaTheme())
	fyne.CurrentApp().SetIcon(resource.WindowIcon)

	//create window
	window := ui.NewMainWindow()

	//create system tray
	if desktop, ok := fyne.CurrentApp().(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
	}

	//load local config
	if err := client.Manager().Initialize(); err != nil {
		logger.Error(err, 0)
		return
	}
	window.ShowAndRun()
}
