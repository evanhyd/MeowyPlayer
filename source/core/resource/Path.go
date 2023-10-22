package resource

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/json"

	"fyne.io/fyne/v2"
)

const (
	collectionFolder = "collection"
	coverFolder      = "cover"
	musicFolder      = "music"
	assetFolder      = "asset"
	collectionFile   = "collection.json"
)

func CollectionPath() string {
	return filepath.Join(collectionFolder, collectionFile)
}

func CoverPath(album *Album) string {
	return filepath.Join(collectionFolder, coverFolder, album.Title+".png")
}

func MusicPath(music *Music) string {
	return filepath.Join(musicFolder, music.Title)
}

func GetCover(album *Album) fyne.Resource {
	asset, err := fyne.LoadResourceFromPath(CoverPath(album))
	if err != nil {
		log.Println(err)
		return MissingIcon
	}
	return asset
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(filepath.Join(collectionFolder, coverFolder), 0777), "failed to create cover folder")
	assert.NoErr(os.MkdirAll(filepath.Join(musicFolder), 0777), "failed to create music folder")

	if _, err := os.Stat(CollectionPath()); os.IsNotExist(err) {
		// create default collection
		assert.NoErr(json.WriteFile(CollectionPath(), &Collection{Date: time.Now(), Albums: make(map[string]Album)}), "failed to create default collection file")
	} else {
		assert.NoErr(err, "failed to fetch collection file info")
	}
}
