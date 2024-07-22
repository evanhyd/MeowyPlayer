package main

import (
	_ "net/http/pprof"

	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view"
)

func main() {
	//curl http://localhost/debug/pprof/heap -O profile.log
	//go tool pprof profile.log
	//go http.ListenAndServe("localhost:80", nil)

	model.InitClient(model.NewLocalStorage())
	player.InitPlayer()
	view.RunApp()
}
