package main

import (
	_ "net/http/pprof"

	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view"

	"fyne.io/fyne/v2"
)

func main() {
	//curl http://localhost/debug/pprof/heap -O profile.log
	//go tool pprof profile.log
	//go http.ListenAndServe("localhost:80", nil)

	if err := model.InitUIClient(); err != nil {
		fyne.LogError("failed to initialize the UI client", err)
		return
	}
	if err := model.InitNetworkClient(); err != nil {
		fyne.LogError("failed to initialize the UI client", err)
		return
	}

	player.InitPlayer()
	view.RunApp()
}
