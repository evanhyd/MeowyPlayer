package custom_canvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SearchList[T any] struct {
	fyne.Container
	searchEntry *widget.Entry   //searching bar
	items       *fyne.Container //displayed items
	DataList    []T             //internal data
}

func NewSearchList[T any](
	placeholder string, //search bar hint
	satisfyQuery func(string, *T) bool, //true if the data satisfies the condition
	makeItem func(*T) fyne.CanvasObject, //constructor of the displayed item
) *SearchList[T] {

	lst := SearchList[T]{}

	//scrolling wrapper
	lst.items = container.NewVBox()
	scroll := container.NewScroll(lst.items)

	//init search bar
	lst.searchEntry = widget.NewEntry()
	lst.searchEntry.SetPlaceHolder(placeholder)

	//add to the displayed item list if fits the query
	lst.searchEntry.OnChanged = func(query string) {

		//clear the old items
		lst.items.Objects = make([]fyne.CanvasObject, 0)

		//check if new data satisfies the query
		for i := 0; i < len(lst.DataList); i++ {
			if satisfyQuery(query, &lst.DataList[i]) {
				lst.items.Add(makeItem(&lst.DataList[i]))
			}
		}
		scroll.ScrollToTop()
	}

	//initialize search list
	lst.Container = *container.NewBorder(lst.searchEntry, nil, nil, nil, scroll)
	lst.DataList = make([]T, 0)
	return &lst
}

/*
   Add data to the data list
*/
func (s *SearchList[T]) AddData(data T) {
	s.DataList = append(s.DataList, data)
}

func (s *SearchList[T]) ClearData() {
	s.DataList = make([]T, 0)
}

func (s *SearchList[T]) ResetSearch() {
	s.searchEntry.SetText("")
	s.searchEntry.OnChanged("")
}
