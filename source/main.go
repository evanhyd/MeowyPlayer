package main

import (
	"encoding/json"
	"log/slog"
	"meowyplayer/loggers"
	"meowyplayer/mcontext"
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

	// Config.
	config := mcontext.UserConfig{}
	if configData, err := os.ReadFile(filepath.Join(baseDir, "config.json")); err != nil {
		slog.Error("failed to read config, fallback to default", "error", err)
		config.Endpoints = map[string]string{
			"register":       `http://40.233.108.102/auth/register`,
			"login":          `http://40.233.108.102/auth/login`,
			"refresh":        `http://40.233.108.102/auth/refresh`,
			"reset-password": `http://40.233.108.102/auth/reset-password`,
			"me":             `http://40.233.108.102/users/me`,
		}
	} else if err = json.Unmarshal(configData, &config); err != nil {
		slog.Error("failed to parse config", "error", err)
		return
	}

	// Storage.
	storage := storages.NewSQLiteStorage(filepath.Join(baseDir, "local.db"), filepath.Join(baseDir, "music"))
	userContext := mcontext.MakeUserContext(config, storage)

	// Start the app.
	ui.RunApp(&userContext)
}
