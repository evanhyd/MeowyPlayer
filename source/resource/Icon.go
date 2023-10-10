package resource

import (
	"os"

	"fyne.io/fyne/v2"
	"meowyplayer.com/utility/assert"
)

const (
	iconNameMissing          = "missing_asset.png"
	iconNameWindow           = "icon.ico"
	iconNameAlbumTab         = "album_tab.png"
	iconNameAlbumAdderOnline = "album_adder_online.png"
	iconNameMusicTab         = "music_tab.png"
	iconNameMusicAdderOnline = "music_adder_online.png"
	iconNameDefault          = "default.png"
	iconRandom               = "random.png"
	iconYouTube              = "youtube.png"
	iconBiliBili             = "bilibili.png"
)

func getResource(resourcePath string) fyne.Resource {
	asset, err := fyne.LoadResourceFromPath(resourcePath)
	if os.IsNotExist(err) {
		asset, err = fyne.LoadResourceFromPath(AssetPath(iconNameMissing))
	}
	assert.NoErr(err, "failed to load icon resource")
	return asset
}

func getIcon(iconName string) fyne.Resource {
	return getResource(AssetPath(iconName))
}

func GetCover(album *Album) fyne.Resource {
	return getResource(CoverPath(album))
}

func WindowIcon() fyne.Resource {
	return getIcon(iconNameWindow)
}

func AlbumTabIcon() fyne.Resource {
	return getIcon(iconNameAlbumTab)
}

func AlbumAdderOnlineIcon() fyne.Resource {
	return getIcon(iconNameAlbumAdderOnline)
}

func MusicTabIcon() fyne.Resource {
	return getIcon(iconNameMusicTab)
}

func MusicAdderOnlineIcon() fyne.Resource {
	return getIcon(iconNameMusicAdderOnline)
}

func DefaultIcon() fyne.Resource {
	return getIcon(iconNameDefault)
}

func RandomIcon() fyne.Resource {
	return getIcon(iconRandom)
}

func YouTubeIcon() fyne.Resource {
	return getIcon(iconYouTube)
}

func BiliBiliIcon() fyne.Resource {
	return getIcon(iconBiliBili)
}
