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
	searchBar.ExtendBaseWidget(searchBar)
	return searchBar
}

func (searchBar *SearchBar) TypedRune(r rune) {
	searchBar.Entry.TypedRune(r)
	searchBar.NotifyAll(searchBar.Text)
}
