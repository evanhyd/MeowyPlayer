package main

import (
	"log/slog"
	"meowyplayer/loggers"
	"meowyplayer/storages"
	"meowyplayer/ui"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"fyne.io/fyne/v2/app"
)

func main() {
	// Panic handler.
	defer func() {
		if e := recover(); e != nil {
			slog.Error("fyne crashed",
				"error", e,
				"stack", string(debug.Stack()),
			)
		}
	}()

	meowApp := app.New()

	// Get the writable sandbox path.
	baseDir := meowApp.Storage().RootURI().Path()
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		slog.Error("failed to create meowyplayer base directory", "error", err)
		return
	}

	// Log path.
	logFilePath := filepath.Join(baseDir, "log.txt")
	logger := loggers.InitializeGlobalLogger(logFilePath)
	defer logger.Close()

	// User config.
	config := storages.GetUserConfig(baseDir)

	// Network client with low timeout limit.
	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{Timeout: 200 * time.Millisecond}).DialContext,
		},
		Timeout: 0,
	}

	userContext := storages.MakeUserContext(
		&httpClient,
		config,
		storages.NewServerStorage(
			&httpClient,
			storages.NewSQLiteStorage(filepath.Join(baseDir, "local.db"), filepath.Join(baseDir, "music")),
			&config,
		),
	)

	ui.RunApp(meowApp, &userContext)
}
