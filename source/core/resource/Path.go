package resource

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"meowyplayer.com/utility/logger"
)

const (
	musicPath = "music"

	collectionPath = "collection"
	collectionFile = "collection.json"
	coverPath      = "cover"

	configFile = "config.json"
)

func SanatizeFileName(filename string) string {
	//sanitize music title
	sanitizer := strings.NewReplacer(
		"<", "",
		">", "",
		":", "",
		"\"", "",
		"/", "",
		"\\", "",
		"|", "",
		"?", "",
		"*", "",
		"~", "",
	)
	return sanitizer.Replace(filename)
}

func MusicPath(music *Music) string {
	return filepath.Join(musicPath, music.Title)
}

func CollectionPath() string {
	return collectionPath
}

func CollectionFile() string {
	return filepath.Join(collectionPath, collectionFile)
}

func CoverPath(album *Album) string {
	return filepath.Join(collectionPath, coverPath, album.Title+".png")
}

func Cover(album *Album) fyne.Resource {
	asset, err := fyne.LoadResourceFromPath(CoverPath(album))
	if err != nil {
		logger.Error(err, 1)
		return MissingIcon
	}
	return asset
}

func ConfigFile() string {
	return configFile
}

func MakeNecessaryPath() {
	if err := os.MkdirAll(filepath.Join(collectionPath, coverPath), 0777); err != nil {
		logger.Error(err, 0)
	}

	if err := os.MkdirAll(filepath.Join(musicPath), 0777); err != nil {
		logger.Error(err, 0)
	}
}
