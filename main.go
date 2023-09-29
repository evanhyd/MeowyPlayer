package main

import (
	"log"
	"runtime/debug"

	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui"
	"meowyplayer.com/source/utility"
)

func main() {
	//redirect panic message
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v\n%v", err, string(debug.Stack()))
		}
	}()

	utility.InitLogger()
	resource.MakeNecessaryPath()

	window := ui.NewMainWindow()
	inUse, err := client.LoadFromLocalCollection()
	utility.MustNil(err)
	client.GetCurrentCollection().Set(&inUse)

	window.ShowAndRun()
}
