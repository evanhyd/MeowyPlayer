package resource

import (
	"log"
	"os"
	"path/filepath"
)

const (
	basePath         = ""
	resourceRootPath = "resource"

	albumRootPath   = "album"
	albumIconName   = "icon.png"
	albumConfigName = "config.txt"

	musicRootPath = "music"
)

func init() {
	if err := os.Mkdir(resourceRootPath, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err := os.Mkdir(albumRootPath, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err := os.Mkdir(musicRootPath, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func GetResourcePath(resource string) string {
	return filepath.Join(basePath, resourceRootPath, resource)
}

func GetAlbumRootPath() string {
	return filepath.Join(basePath, albumRootPath)
}

func GetAlbumFolderPath(album string) string {
	return filepath.Join(GetAlbumRootPath(), album)
}

func GetAlbumIconPath(album string) string {
	return filepath.Join(GetAlbumRootPath(), album, albumIconName)
}

func GetAlbumConfigPath(album string) string {
	return filepath.Join(GetAlbumRootPath(), album, albumConfigName)
}

func GetMusicRootPath() string {
	return filepath.Join(basePath, musicRootPath)
}

func GetMusicPath(music string) string {
	return filepath.Join(GetMusicRootPath(), music)
}
