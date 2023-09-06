package cbinding

import (
	"fyne.io/fyne/v2/data/binding"
	"meowyplayer.com/source/player"
)

type MusicList struct {
	viewBase[player.Music]
}

func NewMusicList() *MusicList {
	return &MusicList{
		viewBase[player.Music]{
			binding.NewUntypedList(),
			nil,
			func(player.Music) bool { return true },
			func(t1, t2 player.Music) bool { return true }},
	}
}

func (c *MusicList) Notify(album *player.Album) {
	c.data = album.MusicList
	c.updateBinding()
}
