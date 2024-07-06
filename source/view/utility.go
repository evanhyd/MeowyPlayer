package view

import (
	"slices"
)

type Titler interface {
	Title() string
}

type dataPipeline[T Titler] struct {
	comparator func(T, T) int
	filter     func(string) bool
}

func newDataPipeline[T Titler]() dataPipeline[T] {
	return dataPipeline[T]{
		comparator: func(T, T) int { return -1 },
		filter:     func(string) bool { return true },
	}
}

func (p *dataPipeline[T]) sortCopy(data []T) []T {
	copy := slices.Clone(data)
	slices.SortStableFunc(copy, p.comparator)
	return copy
}

func (p *dataPipeline[T]) apply(data []T) []T {
	//lazy removal, swapping to the back
	i, j := 0, len(data)
	for i < j {
		if p.filter(data[i].Title()) {
			i++
		} else {
			j--
			data[i], data[j] = data[j], data[i]
		}
	}

	slices.SortStableFunc(data[:j], p.comparator)
	return data[:j]
}
