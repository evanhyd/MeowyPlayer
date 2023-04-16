package resource

import "fyne.io/fyne/v2"

var mainWindowSize fyne.Size
var albumCoverIconSize fyne.Size

func init() {
	mainWindowSize = fyne.NewSize(460, 650)
	albumCoverIconSize = fyne.NewSize(128.0, 128.0)
}

func GetMainWindowSize() fyne.Size {
	return mainWindowSize
}

func GetAlbumCoverSize() fyne.Size {
	return albumCoverIconSize
}
