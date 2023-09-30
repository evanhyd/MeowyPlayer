package resource

import (
	"os"
	"path/filepath"
	"time"

	"meowyplayer.com/source/player"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/json"
)

const (
	albumPath      = "album"
	coverPath      = "cover"
	collectionFile = "collection.json"

	musicPath = "music"
	assetPath = "asset"
)

func CollectionPath() string {
	return filepath.Join(albumPath, collectionFile)
}

func CoverPath(album *player.Album) string {
	return filepath.Join(albumPath, coverPath, album.Title+".png")
}

func MusicPath(music *player.Music) string {
	return filepath.Join(musicPath, music.Title)
}

func AssetPath(assetName string) string {
	return filepath.Join(assetPath, assetName)
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(filepath.Join(albumPath, coverPath), os.ModePerm))
	assert.NoErr(os.MkdirAll(filepath.Join(musicPath), os.ModePerm))

	_, err := os.Stat(CollectionPath())
	if os.IsNotExist(err) {
		//create default collection
		assert.NoErr(json.Write(CollectionPath(), &player.Collection{Date: time.Now(), Albums: nil}))
	} else {
		assert.NoErr(err)
	}
}
