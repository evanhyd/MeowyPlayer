package client

import (
	"slices"

	"meowyplayer.com/source/resource"
	"meowyplayer.com/utility/json"
	"meowyplayer.com/utility/pattern"
)

var collectionData pattern.Data[*resource.Collection]
var albumData pattern.Data[*resource.Album]
var playListData pattern.Data[*resource.PlayList]

// the album pointer parameter may refer to a temporary object from the view list
// we need the original one from the collection
func getSourceAlbum(album *resource.Album) *resource.Album {
	index := slices.IndexFunc(collectionData.Get().Albums, func(a resource.Album) bool { return a.Title == album.Title })
	return &collectionData.Get().Albums[index]
}

func reloadCollectionData() error {
	if err := json.WriteFile(resource.CollectionPath(), collectionData.Get()); err != nil {
		return err
	}
	collection, err := LoadFromLocalCollection()
	if err != nil {
		return err
	}
	collectionData.Set(&collection)
	return nil
}

func reloadAlbumData() error {
	albumData.Set(getSourceAlbum(albumData.Get()))
	return nil
}

func GetCollectionData() *pattern.Data[*resource.Collection] {
	return &collectionData
}

func GetAlbumData() *pattern.Data[*resource.Album] {
	return &albumData
}

func GetPlayListData() *pattern.Data[*resource.PlayList] {
	return &playListData
}

func LoadFromLocalCollection() (resource.Collection, error) {
	inUse := resource.Collection{}
	if err := json.ReadFile(resource.CollectionPath(), &inUse); err != nil {
		return inUse, err
	}

	for i := range inUse.Albums {
		inUse.Albums[i].Cover = resource.GetCover(&inUse.Albums[i])
	}

	return inUse, nil
}
