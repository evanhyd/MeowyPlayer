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
	albumFolderPath = "album"
	coverFolderPath = "cover"
	musicFolderPath = "music"
	assetFolderPath = "asset"

	collectionFile = "collection.json"
)

func CollectionPath() string {
	return filepath.Join(albumFolderPath, collectionFile)
}

func CoverPath(album *Album) string {
	return filepath.Join(albumFolderPath, coverFolderPath, album.Title+".png")
}

func MusicPath(music *Music) string {
	return filepath.Join(musicFolderPath, music.Title)
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
	assert.NoErr(os.MkdirAll(filepath.Join(albumFolderPath, coverFolderPath), os.ModePerm), "failed to create cover folder")
	assert.NoErr(os.MkdirAll(filepath.Join(musicFolderPath), os.ModePerm), "failed to create music folder")

	if _, err := os.Stat(CollectionPath()); os.IsNotExist(err) {
		//create default collection
		assert.NoErr(json.WriteFile(CollectionPath(), &Collection{Date: time.Now(), Albums: nil}), "failed to create default collection file")
	} else {
		assert.NoErr(err, "failed to fetch collection file info")
	}
}
