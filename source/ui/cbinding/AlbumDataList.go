package cbinding

import "meowyplayer.com/source/resource"

type AlbumDataList struct {
	dataList[resource.Album]
}

func MakeAlbumDataList() AlbumDataList {
	return AlbumDataList{makeDataList[resource.Album]()}
}

func (a *AlbumDataList) Notify(collection *resource.Collection) {
	a.dataList.Notify(collection.Albums)
}
