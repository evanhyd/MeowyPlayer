package main

import (
	"log"
	"runtime/debug"

	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/logger"
)

func main() {
	//redirect panic message
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v\n%v", err, string(debug.Stack()))
		}
	}()

	logger.Initiate()
	resource.MakeNecessaryPath()

	window := ui.NewMainWindow()
	inUse, err := client.LoadFromLocalCollection()
	assert.NoErr(err)
	client.GetCollectionData().Set(&inUse)

	window.ShowAndRun()
}
