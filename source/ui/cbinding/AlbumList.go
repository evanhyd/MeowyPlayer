package cbinding

import (
	"fyne.io/fyne/v2/data/binding"
	"meowyplayer.com/source/player"
)

type AlbumList struct {
	viewBase[player.Album]
}

func NewAlbumList() *AlbumList {
	return &AlbumList{
		viewBase[player.Album]{
			binding.NewUntypedList(),
			nil,
			func(player.Album) bool { return true },
			func(t1, t2 player.Album) bool { return true }},
	}
}

func (c *AlbumList) Notify(config *player.Config) {
	c.data = config.Albums
	c.updateBinding()
}
