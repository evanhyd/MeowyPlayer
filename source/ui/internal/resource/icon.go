package resource

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed asset/icon.png
var windowIcon []byte
var windowIconResource = fyne.StaticResource{StaticName: "icon.png", StaticContent: windowIcon}

func WindowIcon() fyne.Resource {
	return &windowIconResource
}

//go:embed asset/collection_tab.svg
var collectionTabIcon []byte
var collectionTabIconResource = fyne.StaticResource{StaticName: "collection_tab.svg", StaticContent: collectionTabIcon}

func CollectionTabIcon() fyne.Resource {
	return &collectionTabIconResource
}

//go:embed asset/youtube.svg
var youtubeIcon []byte
var youtubeIconResource = fyne.StaticResource{StaticName: "youtube.svg", StaticContent: youtubeIcon}

func YouTubeIcon() fyne.Resource {
	return &youtubeIconResource
}

//go:embed asset/alphabetical.svg
var alphabeticalIcon []byte
var alphabeticalIconResource = fyne.StaticResource{StaticName: "alphabetical.svg", StaticContent: alphabeticalIcon}

func AlphabeticalIcon() fyne.Resource {
	return &alphabeticalIconResource
}

//go:embed asset/random.svg
var randomIcon []byte
var randomIconResource = fyne.StaticResource{StaticName: "random.svg", StaticContent: randomIcon}

func RandomIcon() fyne.Resource {
	return &randomIconResource
}
