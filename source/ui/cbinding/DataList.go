package cbinding

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/utility"
)

type DataList[T any] struct {
	utility.Subject[[]T]
	data   []T
	filter func(T) bool
	sorter func(t1, t2 T) bool
}

func MakeDataList[T any]() DataList[T] {
	return DataList[T]{filter: func(t T) bool { return true }, sorter: func(t1, t2 T) bool { return false }}
}

func (d *DataList[T]) SetFilter(filter func(T) bool) {
	d.filter = filter
	d.updateBinding()
}

func (d *DataList[T]) SetSorter(sorter func(t1, t2 T) bool) {
	d.sorter = sorter
	d.updateBinding()
}

func (d *DataList[T]) Notify(data []T) {
	d.data = data
	d.updateBinding()
}

func (d *DataList[T]) updateBinding() {
	slices.SortStableFunc(d.data, d.sorter)

	views := []T{}
	for i := range d.data {
		if d.filter(d.data[i]) {
			views = append(views, d.data[i])
		}
	}
	d.Subject.NotifyAll(views)
}
