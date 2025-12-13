package main

import (
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/events"
	"meowyplayer/frontend"
	"meowyplayer/loggers"
	"meowyplayer/players"
	"meowyplayer/storages"
	"os"
	"path/filepath"
)

func main() {
	// Base directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("failed to get user home directory", "error", err)
		return
	}
	baseDir := filepath.Join(homeDir, "meowyplayer")
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		slog.Error("failed to create meowyplayer base directory", "error", err)
		return
	}

	// Log path.
	logFilePath := filepath.Join(baseDir, "log.txt")
	logger := loggers.InitializeGlobalLogger(logFilePath)
	defer logger.Close()

	// Local storage path.
	dbPath := filepath.Join(baseDir, "local.db")
	musicFilePath := filepath.Join(baseDir, "music")

	player := players.MakeBeepPlayer()
	userContext := context.MakeUserContext()
	userContext.AddListener(events.StorageSetEvent, player.HandleStorageSetEvent)
	userContext.AddListener(events.PlaylistSetEvent, player.HandlePlaylistSetEvent)
	userContext.SetStorage(storages.NewSQLiteStorage(dbPath, musicFilePath, storages.User{UserID: 0}))

	frontend.RunApp()
}
