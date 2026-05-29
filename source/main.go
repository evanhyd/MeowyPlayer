package main

import (
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/loggers"
	"meowyplayer/storages"
	"meowyplayer/ui"
	"os"
	"path/filepath"
)

func main() {
	// Base path.
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
	userContext := context.MakeUserContext()

	ui.RunApp(&userContext, func() {
		userContext.SetStorage(storages.NewSQLiteStorage(dbPath, musicFilePath))
	})
}
