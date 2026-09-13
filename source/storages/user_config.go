package storages

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type EndpointProvider interface {
	GetEndpoint(key string) (string, error)
}

type UserConfig struct {
	Endpoints map[string]string `json:"endpoints"`
}

func (u *UserConfig) GetEndpoint(key string) (string, error) {
	value, ok := u.Endpoints[key]
	if ok && value != "" {
		return value, nil
	}
	err := fmt.Errorf("endpoint %v is not configured", key)
	slog.Error(err.Error())
	return "", err
}

func GetUserConfig(baseDir string) UserConfig {
	config := UserConfig{}

	// Endpoints.
	if configData, err := os.ReadFile(filepath.Join(baseDir, "config.json")); err != nil {
		// Fallback to default if can't find the config file.
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
			"putMusicBulk":            `http://147.5.96.83/api/putMusicBulk`,
			"getPlaylistsFromUser":    `http://147.5.96.83/api/getPlaylistsFromUser`,
			"putMusicInPlaylist":      `http://147.5.96.83/api/putMusicInPlaylist`,
			"putMusicInPlaylistBulk":  `http://147.5.96.83/api/putMusicInPlaylistBulk`,
			"deleteMusicFromPlaylist": `http://147.5.96.83/api/deleteMusicFromPlaylist`,
		}
	} else if err = json.Unmarshal(configData, &config); err != nil {
		// Error if the config is corrupted.
		slog.Error("failed to parse config", "error", err)
	}

	return config
}
