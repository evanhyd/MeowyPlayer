package view

import (
	"slices"

	"fyne.io/fyne/v2"
)

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

type Titler interface {
	Title() string
}

type DataPipeline[T Titler] struct {
	comparator func(T, T) int
	filter     func(string) bool
}

func NewDataPipeline[T Titler]() DataPipeline[T] {
	return DataPipeline[T]{
		comparator: func(T, T) int { return -1 },
		filter:     func(string) bool { return true },
	}
}

func (p *DataPipeline[T]) Pass(data []T) []T {
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
