package main

import (
	// _ "net/http/pprof"

	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view"

	"fyne.io/fyne/v2"
)

func main() {
	//go http.ListenAndServe("localhost:80", nil)
	//curl http://localhost/debug/pprof/heap -O profile.log
	//go tool pprof profile.log
	if err := model.InitStorageClient(); err != nil {
		fyne.LogError("failed to initialize the UI client", err)
		return
	}
	if err := model.InitNetworkClient(); err != nil {
		fyne.LogError("failed to initialize the UI client", err)
		return
	}
	if err := player.InitPlayer(); err != nil {
		fyne.LogError("failed to initialize the Player", err)
		return
	}

	view.RunApp()
}
