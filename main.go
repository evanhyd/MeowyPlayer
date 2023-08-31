package main

import (
	"log"
	"os"
	"runtime/debug"

	"meowyplayer.com/source/resource"
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

	window := ui.NewMainWindow()
	if config, err := resource.LoadFromLocalConfig(); err == nil || os.IsNotExist(err) {
		resource.SetCurrentConfig(&config)
	} else {
		log.Panic(err)
	}

	window.ShowAndRun()
}
