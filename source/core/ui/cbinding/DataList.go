package cbinding

import (
	"slices"

	"meowyplayer.com/utility/container"
	"meowyplayer.com/utility/pattern"
)

type dataList[T any] struct {
	pattern.SubjectBase[[]T]
	data   container.Slice[T]
	filter func(T) bool
	sorter func(T, T) int
}

func makeDataList[T any]() dataList[T] {
	return dataList[T]{filter: func(T) bool { return true }, sorter: func(T, T) int { return -1 }}
}

func (d *dataList[T]) SetFilter(filter func(T) bool) {
	d.filter = filter
	d.updateBinding()
}

func (d *dataList[T]) SetSorter(sorter func(T, T) int) {
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
