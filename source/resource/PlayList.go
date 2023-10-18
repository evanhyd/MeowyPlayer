package resource

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/utility/assert"
)

type PlayList struct {
	musicList []Music
	index     int
}

func NewPlayList(musicList []Music, music *Music) *PlayList {
	index := slices.Index(musicList, *music)
	return &PlayList{musicList, index}
}

func (p *PlayList) Music() *Music {
	return &p.musicList[p.index]
}

func (p *PlayList) MusicCount() int {
	return len(p.musicList)
}

func (p *PlayList) Index() int {
	return p.index
}

func (p *PlayList) SetIndex(musicIndex int) {
	assert.Ensure(func() bool { return 0 <= musicIndex && musicIndex < len(p.musicList) })
	p.index = musicIndex
}
