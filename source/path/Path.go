package path

import (
	"path/filepath"

	"meowyplayer.com/source/player"
)

const (
	albumPath  = "album"
	coverPath  = "cover"
	configFile = "config.json"

	musicPath = "music"

	assetPath = "asset"
)

func Config() string {
	return filepath.Join(albumPath, configFile)
}

func Cover(album *player.Album) string {
	return filepath.Join(albumPath, coverPath, album.Title+".png")
}

func Music(music *player.Music) string {
	return filepath.Join(musicPath, music.Title)
}

func Asset(assetName string) string {
	return filepath.Join(assetPath, assetName)
}
