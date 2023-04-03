package cwidget

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/pattern"
)

type List[T any] struct {
	widget.List
	internalData    []T
	displayedData   []T
	onTappedSubject pattern.OneArgSubject[T]
	filter          func(T) bool
	sorter          func(T, T) bool
}

func NewList[T any](createItem func() fyne.CanvasObject, updateItem func(T, fyne.CanvasObject)) *List[T] {
	list := &List[T]{}
	list.Initialize(createItem, updateItem)
	list.ExtendBaseWidget(list)
	return list
}

func (list *List[T]) Initialize(createItem func() fyne.CanvasObject, updateItem func(T, fyne.CanvasObject)) {
	list.Length = func() int { return len(list.displayedData) }
	list.CreateItem = createItem
	list.UpdateItem = func(id widget.ListItemID, canvas fyne.CanvasObject) { updateItem(list.displayedData[id], canvas) }

	//primary tap
	list.List.OnSelected = func(id widget.ListItemID) {
		list.onTappedSubject.NotifyAll(list.displayedData[id])
		list.Unselect(id)
	}

	//secondary tap

	list.filter = func(T) bool { return true }
	list.sorter = func(T, T) bool { return true }
}

func (list *List[T]) Notify(items []T) {
	list.internalData = items
	list.refreshDisplayData()
	list.ScrollToTop()
}

func (list *List[T]) OnTappedSubject() *pattern.OneArgSubject[T] {
	return &list.onTappedSubject
}

func (list *List[T]) SetOnTapped(onTapped func(T)) {
	list.List.OnSelected = func(id widget.ListItemID) {
		onTapped(list.displayedData[id])
		list.onTappedSubject.NotifyAll(list.displayedData[id])
		list.UnselectAll()
	}
}

func (list *List[T]) SetFilter(filter func(T) bool) {
	list.filter = filter
	list.refreshDisplayData()
}

func (list *List[T]) SetSorter(sorter func(T, T) bool) {
	list.sorter = sorter
	list.refreshDisplayData()
}

func (list *List[T]) refreshDisplayData() {
	list.displayedData = list.displayedData[0:0]

	//sort first, since internal data may affect the playing order
	sort.SliceStable(list.internalData, func(i, j int) bool {
		return list.sorter(list.internalData[i], list.internalData[j])
	})

	//filter keeps satisfied data
	for i := range list.internalData {
		if list.filter(list.internalData[i]) {
			list.displayedData = append(list.displayedData, list.internalData[i])
		}
	}

	list.Refresh()
}
