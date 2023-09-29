package player

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/utility"
)

type PlayList struct {
	album Album
	index int
}

func NewPlayList(album *Album, music *Music) *PlayList {
	index := slices.Index(album.MusicList, *music)
	utility.Assert(func() bool { return index != -1 })
	return &PlayList{*album, index}
}

func (p *PlayList) Album() *Album {
	return &p.album
}

func (p *PlayList) Music() *Music {
	return &p.album.MusicList[p.index]
}

func (p *PlayList) Index() int {
	return p.index
}

func (p *PlayList) SetIndex(musicIndex int) {
	utility.Assert(func() bool { return 0 <= musicIndex && musicIndex < len(p.album.MusicList) })
	p.index = musicIndex
}
