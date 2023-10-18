package cbinding

import (
	"golang.org/x/exp/maps"
	"meowyplayer.com/source/resource"
)

type AlbumDataList struct {
	dataList[resource.Album]
}

func MakeAlbumDataList() AlbumDataList {
	return AlbumDataList{makeDataList[resource.Album]()}
}

func (a *AlbumDataList) Notify(collection resource.Collection) {
	a.dataList.Notify(maps.Values(collection.Albums))
}
