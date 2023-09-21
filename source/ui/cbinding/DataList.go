package cbinding

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/utility"
)

type dataList[T any] struct {
	utility.Subject[[]T]
	data   []T
	filter func(T) bool
	sorter func(t1, t2 T) bool
}

func makeDataList[T any]() dataList[T] {
	return dataList[T]{filter: func(t T) bool { return true }, sorter: func(t1, t2 T) bool { return false }}
}

func (d *dataList[T]) SetFilter(filter func(T) bool) {
	d.filter = filter
	d.updateBinding()
}

func (d *dataList[T]) SetSorter(sorter func(t1, t2 T) bool) {
	d.sorter = sorter
	d.updateBinding()
}

func (d *dataList[T]) Notify(data []T) {
	d.data = data
	d.updateBinding()
}

func (d *dataList[T]) updateBinding() {
	slices.SortStableFunc(d.data, d.sorter)

	views := []T{}
	for i := range d.data {
		if d.filter(d.data[i]) {
			views = append(views, d.data[i])
		}
	}
	d.Subject.NotifyAll(views)
}
