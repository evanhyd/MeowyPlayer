package cbinding

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/utility/container"
	"meowyplayer.com/utility/pattern"
)

type dataList[T any] struct {
	pattern.Subject[[]T]
	data   container.Slice[T]
	filter func(T) bool
	sorter func(T, T) bool
}

func makeDataList[T any]() dataList[T] {
	return dataList[T]{filter: func(T) bool { return true }, sorter: func(T, T) bool { return false }}
}

func (d *dataList[T]) SetFilter(filter func(T) bool) {
	d.filter = filter
	d.updateBinding()
}

func (d *dataList[T]) SetSorter(sorter func(T, T) bool) {
	d.sorter = sorter
	d.updateBinding()
}

func (d *dataList[T]) Notify(data []T) {
	d.data = data
	d.updateBinding()
}

func (d *dataList[T]) updateBinding() {
	slices.SortStableFunc(d.data, d.sorter)
	d.NotifyAll(d.data.Filter(d.filter))
}
