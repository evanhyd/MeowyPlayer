package cwidget

import (
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/pattern"
)

type SearchBar struct {
	widget.Entry
	pattern.OneArgSubject[string]
}

func NewSearchBar() *SearchBar {
	searchBar := &SearchBar{}
	searchBar.OnChanged = searchBar.NotifyAll
	searchBar.ExtendBaseWidget(searchBar)
	return searchBar
}

func (searchBar *SearchBar) SetOnChanged(onChanged func(string)) {
	searchBar.OnChanged = func(text string) {
		onChanged(text)
		searchBar.NotifyAll(text)
	}
}
