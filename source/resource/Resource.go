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

func GetAsset(assetName string) fyne.Resource {
	return get(path.Asset(assetName))
}

func get(resourcePath string) fyne.Resource {
	const missingAssetName = "missing_asset.png"

	asset, err := fyne.LoadResourceFromPath(resourcePath)
	if os.IsNotExist(err) {
		asset, err = fyne.LoadResourceFromPath(path.Asset(missingAssetName))
	}
	utility.MustNil(err)

	return asset
}
