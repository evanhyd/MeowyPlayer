package cwidget

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/pattern"
)

type ItemList[T any] struct {
	widget.List
	internalData  []T
	displayedData []T
	onSelected    pattern.OneArgSubject[T]
	filter        func(T) bool
	sorter        func(T, T) bool
}

func NewItemList[T any](createItem func() fyne.CanvasObject, updateItem func(T, fyne.CanvasObject)) *ItemList[T] {
	itemList := &ItemList[T]{}
	itemList.Initialize(createItem, updateItem)
	itemList.ExtendBaseWidget(itemList)
	return itemList
}

func (itemList *ItemList[T]) Initialize(createItem func() fyne.CanvasObject, updateItem func(T, fyne.CanvasObject)) {
	itemList.Length = func() int { return len(itemList.displayedData) }
	itemList.CreateItem = createItem
	itemList.UpdateItem = func(id widget.ListItemID, canvas fyne.CanvasObject) { updateItem(itemList.displayedData[id], canvas) }
	itemList.List.OnSelected = func(id widget.ListItemID) {
		itemList.onSelected.NotifyAll(itemList.displayedData[id])
		itemList.UnselectAll()
	}
	itemList.filter = func(T) bool { return true }
	itemList.sorter = func(T, T) bool { return true }
}

func (itemList *ItemList[T]) OnSelected() *pattern.OneArgSubject[T] {
	return &itemList.onSelected
}

func (itemList *ItemList[T]) SetOnSelected(onSelected func(T)) {
	itemList.List.OnSelected = func(id widget.ListItemID) {
		onSelected(itemList.displayedData[id])
		itemList.onSelected.NotifyAll(itemList.displayedData[id])
		itemList.UnselectAll()
	}
}

func (itemList *ItemList[T]) Notify(items []T) {
	itemList.internalData = items
	itemList.refreshDisplayData()
	itemList.ScrollToTop()
}

func (itemList *ItemList[T]) SetFilter(filter func(T) bool) {
	itemList.filter = filter
	itemList.refreshDisplayData()
}

func (itemList *ItemList[T]) SetSorter(sorter func(T, T) bool) {
	itemList.sorter = sorter
	itemList.refreshDisplayData()
}

func (itemList *ItemList[T]) refreshDisplayData() {
	itemList.displayedData = itemList.displayedData[0:0]

	//sort first, since internal data may affect the playing order
	sort.SliceStable(itemList.internalData, func(i, j int) bool {
		return itemList.sorter(itemList.internalData[i], itemList.internalData[j])
	})

	//filter keeps satisfied data
	for i := range itemList.internalData {
		if itemList.filter(itemList.internalData[i]) {
			itemList.displayedData = append(itemList.displayedData, itemList.internalData[i])
		}
	}

	itemList.Refresh()
}
