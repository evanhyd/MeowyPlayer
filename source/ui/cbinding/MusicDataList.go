package cbinding

import "meowyplayer.com/source/player"

type MusicDataList struct {
	dataList[player.Music]
	album player.Album
}

func MakeMusicDataList() MusicDataList {
	return MusicDataList{makeDataList[player.Music](), player.Album{}}
}

func (m *MusicDataList) Notify(album *player.Album) {
	m.album = *album
	m.dataList.Notify(album.MusicList)
}

func (m *MusicDataList) GetAlbum() *player.Album {
	return &m.album
}
