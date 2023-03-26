package resource

import (
	"log"
	"os"
	"path/filepath"
)

const (
	basePath     = ""
	resourcePath = "resource"
	albumPath    = "album"
	musicPath    = "music"

	albumIconName   = "icon.png"
	albumConfigName = "config.txt"
)

func init() {
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		if err = os.Mkdir(resourcePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat(albumPath); os.IsNotExist(err) {
		if err = os.Mkdir(albumPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat(musicPath); os.IsNotExist(err) {
		if err = os.Mkdir(musicPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func GetResourcePath(resource string) string {
	return filepath.Join(basePath, resourcePath, resource)
}

func GetAlbumPath() string {
	return filepath.Join(basePath, albumPath)
}

func GetMusicPath() string {
	return filepath.Join(basePath, musicPath)
}

func GetAlbumIconPath(album string) string {
	return filepath.Join(GetAlbumPath(), album, albumIconName)
}

func GetAlbumConfigPath(album string) string {
	return filepath.Join(GetAlbumPath(), album, albumConfigName)
}
