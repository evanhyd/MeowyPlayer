package resource

import (
	"os"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/path"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
)

func GetCover(album *player.Album) fyne.Resource {
	return get(path.Cover(album))
}

func get(resourcePath string) fyne.Resource {
	asset, err := fyne.LoadResourceFromPath(resourcePath)
	if os.IsNotExist(err) {
		asset, err = fyne.LoadResourceFromPath(path.Asset(iconNameMissing))
	}
	utility.MustNil(err)
	return asset
}
