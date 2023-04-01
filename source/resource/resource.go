package resource

import (
	"log"
	"os"
	"path/filepath"
)

const (
	basePath           = ""
	resourceFolderPath = "resource"

	albumFolderPath = "album"
	albumIconName   = "icon.png"
	albumConfigName = "config.txt"

	musicFolderPath = "music"
)

func init() {
	if _, err := os.Stat(resourceFolderPath); os.IsNotExist(err) {
		if err = os.Mkdir(resourceFolderPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat(albumFolderPath); os.IsNotExist(err) {
		if err = os.Mkdir(albumFolderPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat(musicFolderPath); os.IsNotExist(err) {
		if err = os.Mkdir(musicFolderPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func GetResourcePath(resource string) string {
	return filepath.Join(basePath, resourceFolderPath, resource)
}

func GetAlbumFolderPath() string {
	return filepath.Join(basePath, albumFolderPath)
}

func GetAlbumIconPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(), album, albumIconName)
}

func GetAlbumConfigPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(), album, albumConfigName)
}

func GetMusicFolderPath() string {
	return filepath.Join(basePath, musicFolderPath)
}

func GetMusicPath(music string) string {
	return filepath.Join(GetMusicFolderPath(), music)
}
