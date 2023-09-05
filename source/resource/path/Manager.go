package path

import (
	"path/filepath"

	"meowyplayer.com/source/player"
)

const (
	albumPath          = "album"
	iconPath           = "icon"
	configFile         = "config.json"
	texturePath        = "texture"
	missingTexturePath = "missing_texture.png"
)

func Config() string {
	return filepath.Join(albumPath, configFile)
}

func Icon(album *player.Album) string {
	return filepath.Join(albumPath, iconPath, album.Title+".png")
}
