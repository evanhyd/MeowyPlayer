package cbinding

import (
	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
)

type ConfigList struct {
	binding.UntypedList
	data   []player.Album
	filter func(*player.Album) bool
	sorter func(player.Album, player.Album) bool
}

func NewConfigList() *ConfigList {
	return &ConfigList{
		binding.NewUntypedList(),
		nil,
		func(*player.Album) bool { return true },
		func(player.Album, player.Album) bool { return true },
	}
}

func (c *ConfigList) Notify(config *player.Config) {
	c.data = config.Albums
	c.updateBinding()
}

func (c *ConfigList) SetFilter(filter func(*player.Album) bool) {
	c.filter = filter
	c.updateBinding()
}

func (c *ConfigList) SetSorter(sorter func(player.Album, player.Album) bool) {
	c.sorter = sorter
	c.updateBinding()
}

func (c *ConfigList) updateBinding() {
	//stable sort
	slices.SortStableFunc(c.data, c.sorter)

	//filter keeps the wanted album
	view := []any{}
	for i := range c.data {
		if c.filter(&c.data[i]) {
			view = append(view, &c.data[i])
		}
	}
	utility.MustOk(c.Set(nil))
	utility.MustOk(c.Set(view))
}
