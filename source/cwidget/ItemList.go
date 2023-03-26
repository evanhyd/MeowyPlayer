package cwidget

import (
	"log"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ItemList[T any] struct {
	widget.List
	internalData  []T
	displayedData []T
	filter        func(T) bool
	sorter        func(T, T) bool
}

func NewItemList[T any](createItem func() fyne.CanvasObject, updateItem func(T, fyne.CanvasObject)) *ItemList[T] {
	itemList := &ItemList[T]{}
	itemList.Length = func() int { return len(itemList.displayedData) }
	itemList.CreateItem = createItem
	itemList.UpdateItem = func(id widget.ListItemID, canvas fyne.CanvasObject) { updateItem(itemList.displayedData[id], canvas) }
	itemList.filter = func(T) bool { return true }
	itemList.sorter = func(T, T) bool { return true }
	itemList.ExtendBaseWidget(itemList)
	return itemList
}

func (itemList *ItemList[T]) SetItems(items []T) {
	itemList.internalData = items
	itemList.Refresh()
}

func (itemList *ItemList[T]) SetFilter(filter func(T) bool) {
	itemList.filter = filter
	itemList.Refresh()
}

func (itemList *ItemList[T]) SetSorter(sorter func(T, T) bool) {
	itemList.sorter = sorter
	itemList.Refresh()
}

func (itemList *ItemList[T]) Refresh() {
	itemList.displayedData = itemList.displayedData[0:0]

	//filter keeps satisfied data
	for i := range itemList.internalData {
		if itemList.filter(itemList.internalData[i]) {
			itemList.displayedData = append(itemList.displayedData, itemList.internalData[i])
		}
	}

	//sort the display data
	sort.Slice(itemList.displayedData, func(i, j int) bool {
		return itemList.sorter(itemList.displayedData[i], itemList.displayedData[j])
	})

	itemList.List.Refresh()
	log.Println("refresh item list")
}
