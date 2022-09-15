package custom_canvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SearchList[T any] struct {
	fyne.Container
	SearchBar *widget.Entry   //searching bar
	items     *fyne.Container //displayed items
	DataList  []T             //internal data
}

func NewSearchList[T any](
	placeholder string, //search bar
	satisfyQuery func(string, *T) bool, //true if the data satisfies the condition
	makeItem func(*T) fyne.CanvasObject, //constructor of the displayed item
) *SearchList[T] {

	lst := SearchList[T]{}

	//init search bar
	lst.SearchBar = widget.NewEntry()
	lst.SearchBar.SetPlaceHolder(placeholder)

	//add to the displayed item list if fits the query
	lst.SearchBar.OnChanged = func(query string) {

		//clear the old items
		lst.items.Objects = make([]fyne.CanvasObject, 0)

		//check if new data satisfies the query
		for i := 0; i < len(lst.DataList); i++ {
			if satisfyQuery(query, &lst.DataList[i]) {
				lst.items.Add(makeItem(&lst.DataList[i]))
			}
		}
	}

	//initialize search list
	lst.items = container.NewVBox()
	lst.Container = *container.NewBorder(lst.SearchBar, nil, nil, nil, container.NewScroll(lst.items))
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
	s.SearchBar.SetText("")
	s.SearchBar.OnChanged("")
}
