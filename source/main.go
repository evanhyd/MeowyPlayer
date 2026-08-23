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
		slog.Info("missing config file, fallback to default", "info", err)
		config.Endpoints = map[string]string{
			"register":                `http://40.233.108.102/auth/register`,
			"login":                   `http://40.233.108.102/auth/login`,
			"refresh":                 `http://40.233.108.102/auth/refresh`,
			"reset-password":          `http://40.233.108.102/auth/reset-password`,
			"me":                      `http://40.233.108.102/users/me`,
			"getPlaylist":             `http://147.5.96.83/api/getPlaylist`,
			"getPlaylistContent":      `http://147.5.96.83/api/getPlaylistContent`,
			"putPlaylist":             `http://147.5.96.83/api/putPlaylist`,
			"deletePlaylist":          `http://147.5.96.83/api/deletePlaylist`,
			"getMusic":                `http://147.5.96.83/api/getMusic`,
			"getMusicBulk":            `http://147.5.96.83/api/getMusicBulk`,
			"putMusic":                `http://147.5.96.83/api/putMusic`,
			"getPlaylistsFromUser":    `http://147.5.96.83/api/getPlaylistsFromUser`,
			"putMusicInPlaylist":      `http://147.5.96.83/api/putMusicInPlaylist`,
			"putMusicInPlaylistBulk":  `http://147.5.96.83/api/putMusicInPlaylistBulk`,
			"deleteMusicFromPlaylist": `http://147.5.96.83/api/deleteMusicFromPlaylist`,
		}
	} else if err = json.Unmarshal(configData, &config); err != nil {
		slog.Error("failed to parse config", "error", err)
		return
	}

	// Storage.
	localStorage := storages.NewSQLiteStorage(filepath.Join(baseDir, "local.db"), filepath.Join(baseDir, "music"))
	serverStorage := storages.NewServerStorage(localStorage, config.Endpoints)
	userContext := mcontext.MakeUserContext(config, serverStorage)

	// Start the app.
	ui.RunApp(&userContext)
}
