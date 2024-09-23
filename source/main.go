package main

import (
	// _ "net/http/pprof"

	"io"
	"log"
	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view"
	"os"

	"fyne.io/fyne/v2"
)

func main() {
	//go http.ListenAndServe("localhost:80", nil)
	//curl http://localhost/debug/pprof/heap -O profile.log
	//go tool pprof profile.log

	//initialize logger
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fyne.LogError("failed to initialize log file", err)
		return
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	//initialize clients
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
