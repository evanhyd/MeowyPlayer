package path

import (
	"os"
	"path/filepath"
	"time"

	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
)

const (
	albumPath  = "album"
	coverPath  = "cover"
	configFile = "config.json"

	musicPath = "music"
	assetPath = "asset"
)

func Config() string {
	return filepath.Join(albumPath, configFile)
}

func Cover(album *player.Album) string {
	return filepath.Join(albumPath, coverPath, album.Title+".png")
}

func Music(music *player.Music) string {
	return filepath.Join(musicPath, music.Title)
}

func Asset(assetName string) string {
	return filepath.Join(assetPath, assetName)
}

func MakeNecessaryPath() {
	utility.MustNil(os.MkdirAll(filepath.Join(albumPath, coverPath), os.ModePerm))
	utility.MustNil(os.MkdirAll(filepath.Join(musicPath), os.ModePerm))

	_, err := os.Stat(Config())
	if os.IsNotExist(err) {
		utility.MustNil(utility.WriteJson(Config(), &player.Config{Date: time.Now(), Albums: nil}))
	} else {
		utility.MustNil(err)
	}
}
