package cwidget

import "fyne.io/fyne/v2/widget"

type SearchBarObserver interface {
	Notify(text string)
}

type SearchBar struct {
	widget.Entry
	observers []SearchBarObserver
}

func NewSearchBar() *SearchBar {
	searchBar := &SearchBar{}
	searchBar.OnChanged = searchBar.NotifyObservers
	searchBar.ExtendBaseWidget(searchBar)
	return searchBar
}

func (searchBar *SearchBar) SetOnChanged(onChanged func(string)) {
	searchBar.OnChanged = func(text string) {
		onChanged(text)
		searchBar.NotifyObservers(text)
	}
}

func (searchBar *SearchBar) AddObserver(observer SearchBarObserver) {
	searchBar.observers = append(searchBar.observers, observer)
}

func (searchBarthis *SearchBar) NotifyObservers(text string) {
	for _, observer := range searchBarthis.observers {
		observer.Notify(text)
	}
}
