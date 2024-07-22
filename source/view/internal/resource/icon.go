package resource

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed asset/icon.png
var windowIcon []byte
var windowIconResource = fyne.StaticResource{StaticContent: windowIcon}

func WindowIcon() fyne.Resource {
	return &windowIconResource
}

//go:embed asset/collection_tab.svg
var collectionTabIcon []byte
var collectionTabIconResource = fyne.StaticResource{StaticContent: collectionTabIcon}

func CollectionTabIcon() fyne.Resource {
	return &collectionTabIconResource
}

//go:embed asset/youtube.svg
var youtubeIcon []byte
var youtubeIconResource = fyne.StaticResource{StaticContent: youtubeIcon}

func YouTubeIcon() fyne.Resource {
	return &youtubeIconResource
}

//go:embed asset/alphabetical.svg
var alphabeticalIcon []byte
var alphabeticalIconResource = fyne.StaticResource{StaticContent: alphabeticalIcon}

func AlphabeticalIcon() fyne.Resource {
	return &alphabeticalIconResource
}

//go:embed asset/random.svg
var randomIcon []byte
var randomIconResource = fyne.StaticResource{StaticContent: randomIcon}

func RandomIcon() fyne.Resource {
	return &randomIconResource
}
