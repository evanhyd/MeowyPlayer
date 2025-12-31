package ui

import (
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type PlaylistPage struct {
	widget.BaseWidget
	searchEntry *widget.Entry
}

func newPlaylistPage() *PlaylistPage {
	p := PlaylistPage{
		searchEntry: widget.NewEntry(),
	}

	p.searchEntry.SetPlaceHolder(lang.L("Search playlist"))
	p.searchEntry.OnChanged = p.submitSearchQuery
	return &p
}

func (p *PlaylistPage) submitSearchQuery(title string) {

}
