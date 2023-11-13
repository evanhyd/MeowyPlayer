package resource

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"meowyplayer.com/utility/logger"
)

const (
	collectionPath = "collection"
	coverPath      = "cover"
	musicPath      = "music"
	assetPath      = "asset"
	collectionFile = "collection.json"
)

func CollectionPath() string {
	return collectionPath
}

func CollectionFile() string {
	return filepath.Join(collectionPath, collectionFile)
}

func CoverPath(album *Album) string {
	return filepath.Join(collectionPath, coverPath, album.Title+".png")
}

func MusicPath(music *Music) string {
	return filepath.Join(musicPath, music.Title)
}

func Cover(album *Album) fyne.Resource {
	asset, err := fyne.LoadResourceFromPath(CoverPath(album))
	if err != nil {
		logger.Error(err, 1)
		return MissingIcon
	}
	return asset
}

func MakeNecessaryPath() {
	if err := os.MkdirAll(filepath.Join(collectionPath, coverPath), 0777); err != nil {
		logger.Error(err, 0)
	}

	if err := os.MkdirAll(filepath.Join(musicPath), 0777); err != nil {
		logger.Error(err, 0)
	}
}
