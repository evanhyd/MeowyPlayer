package resource

import (
	"fyne.io/fyne/v2"
	"meowyplayer.com/source/path"
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
)

func getIcon(iconName string) fyne.Resource {
	return get(path.Asset(iconName))
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
