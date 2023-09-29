package resource

import (
	"os"
	"path/filepath"
	"time"

	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
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
	utility.MustNil(os.MkdirAll(filepath.Join(albumPath, coverPath), os.ModePerm))
	utility.MustNil(os.MkdirAll(filepath.Join(musicPath), os.ModePerm))

	_, err := os.Stat(CollectionPath())
	if os.IsNotExist(err) {
		utility.MustNil(utility.WriteJson(CollectionPath(), &player.Collection{Date: time.Now(), Albums: nil}))
	} else {
		utility.MustNil(err)
	}
}
