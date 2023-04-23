package resource

import "fyne.io/fyne/v2"

const goldenRatio = 1.61803398875

func GetMainWindowSize() fyne.Size {
	return fyne.NewSize(460.0, 650.0)
}

func GetAlbumCoverSize() fyne.Size {
	return fyne.NewSize(128.0, 128.0)
}

func GetAlbumViewIconSize() fyne.Size {
	return fyne.NewSize(128.0, 128.0)
}

func GetThumbnailIconSize() fyne.Size {
	return fyne.NewSize(128.0*goldenRatio, 128.0)
}

func GetMusicAddOnlineDialogSize() fyne.Size {
	return fyne.NewSize(7680.0, 4320.0)
}
