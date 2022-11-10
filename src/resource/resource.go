package resource

import (
	"os"
	"path/filepath"
	"strings"
)

func GetBasePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if strings.Contains(ex, "var/folders") {
		return ""
	}
	return filepath.Dir(ex)
}

func GetImageBasePath() string {
	return "images"
}

func GetMusicBasePath() string {
	return "music"
}

func GetAlbumBasePath() string {
	return "album"
}

func GetImagePath(image string) string {
	return filepath.Join(GetBasePath(), GetImageBasePath(), image)
}

func GetMusicPath(music string) string {
	return filepath.Join(GetBasePath(), GetMusicBasePath(), music)
}

func GetAlbumFolderPath(album string) string {
	return filepath.Join(GetBasePath(), GetAlbumBasePath(), album)
}

func GetAlbumIconPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(album), "icon.png")
}

func GetAlbumConfigPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(album), "config.txt")
}
