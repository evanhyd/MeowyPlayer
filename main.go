package main

import (
	"fmt"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/logger"
)

func main() {
	// redirect panic message
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Errorf("%v\n%v", err, string(debug.Stack())), "caught panic error", 1)
		}
	}()

	logger.Initiate()
	resource.MakeNecessaryPath()

	//initiate app configuration
	fyne.SetCurrentApp(app.NewWithID("MeowyPlayer"))
	fyne.CurrentApp().Settings().SetTheme(resource.VanillaTheme())

	//create window
	window := ui.NewMainWindow()

	//create system tray
	if desktop, ok := fyne.CurrentApp().(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", window.Show)))
		desktop.SetSystemTrayIcon(resource.WindowIcon)
	}

	//load local config
	collection, err := client.LoadFromLocalCollection()
	assert.NoErr(err, "failed to load from local collection")
	client.GetInstance().SetCollection(collection)
	window.ShowAndRun()
}
