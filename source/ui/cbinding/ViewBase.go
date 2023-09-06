package cbinding

import (
	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/utility"
)

type viewBase[T any] struct {
	binding.UntypedList
	data   []T
	filter func(T) bool
	sorter func(t1, t2 T) bool
}

func (c *viewBase[T]) SetFilter(filter func(T) bool) {
	c.filter = filter
	c.updateBinding()
}

func (c *viewBase[T]) SetSorter(sorter func(t1, t2 T) bool) {
	c.sorter = sorter
	c.updateBinding()
}

func (c *viewBase[T]) updateBinding() {
	slices.SortStableFunc(c.data, c.sorter)

	view := []any{}
	for _, d := range c.data {
		if c.filter(d) {
			view = append(view, d)
		}
	}
	utility.MustOk(c.Set(nil)) //update when changing the length
	utility.MustOk(c.Set(view))
}
