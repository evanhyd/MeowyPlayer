package main

import (
	"log"
	"runtime/debug"

	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/path"
	"meowyplayer.com/source/ui"
	"meowyplayer.com/source/utility"
)

func main() {
	//redirect panic message
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("%v\n%v", err, string(debug.Stack()))
		}
	}()

	utility.InitLogger()
	path.MakeNecessaryPath()

	window := ui.NewMainWindow()
	inUse, err := manager.LoadFromLocalConfig()
	utility.MustNil(err)
	manager.GetCurrentConfig().Set(&inUse)

	window.ShowAndRun()
}
