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
	iconNameAlbumAdderLocal  = "album_adder_local.png"
	iconNameAlbumAdderOnline = "album_adder_online.png"
	iconNameMusicTab         = "music_tab.png"
	iconNameMusicAdderLocal  = "music_adder_local.png"
	iconNameMusicAdderOnline = "music_adder_online.png"
	iconNameSearch           = "search.png"
	iconNameDefault          = "default.png"
	iconPlayModeRandom       = "playmode_random.png"
	iconPlayModeOrdered      = "playmode_ordered.png"
	iconPlayModeRepeat       = "playmode_repeat.png"
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

func AlbumAdderLocalIcon() fyne.Resource {
	return getIcon(iconNameAlbumAdderLocal)
}

func AlbumAdderOnlineIcon() fyne.Resource {
	return getIcon(iconNameAlbumAdderOnline)
}

func MusicTabIcon() fyne.Resource {
	return getIcon(iconNameMusicTab)
}

func MusicAdderLocalIcon() fyne.Resource {
	return getIcon(iconNameMusicAdderLocal)
}

func MusicAdderOnlineIcon() fyne.Resource {
	return getIcon(iconNameMusicAdderOnline)
}

func SearchIcon() fyne.Resource {
	return getIcon(iconNameSearch)
}

func DefaultIcon() fyne.Resource {
	return getIcon(iconNameDefault)
}

func PlayModeRandomIcon() fyne.Resource {
	return getIcon(iconPlayModeRandom)
}

func PlayModeOrderedIcon() fyne.Resource {
	return getIcon(iconPlayModeOrdered)
}

func PlayModeRepeatIcon() fyne.Resource {
	return getIcon(iconPlayModeRepeat)
}
