package main

import (
	"fmt"
	"runtime/debug"

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

	window := ui.NewMainWindow()
	inUse, err := client.LoadFromLocalCollection()
	assert.NoErr(err, "failed to load from local collection")
	client.GetCollectionData().Set(&inUse)
	window.ShowAndRun()
}
