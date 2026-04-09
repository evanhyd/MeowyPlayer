package ui

import (
	"meowyplayer/storages"
	"meowyplayer/ui/internal/layouts"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlaylistPage struct {
	widget.BaseWidget
	searchEntry   *widget.Entry
	searchButton  *widget.Button
	content       *widget.GridWrap
	searchResults []storages.Playlist
}

func newPlaylistPage() *PlaylistPage {
	p := PlaylistPage{
		searchEntry:  widget.NewEntry(),
		searchButton: widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.submitSearchQuery
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.submitSearchQuery(p.searchEntry.Text) }

	p.content = widget.NewGridWrap(
		func() int {
			return 0
		},
		func() fyne.CanvasObject {
			return nil
		},
		func(index widget.GridWrapItemID, object fyne.CanvasObject) {
		},
	)

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *PlaylistPage) CreateRenderer() fyne.WidgetRenderer {
	searchTools := container.New(layouts.NewHSegmentLayout(2, 7, 2), layout.NewSpacer(), p.searchEntry, layout.NewSpacer())
	return widget.NewSimpleRenderer(container.NewBorder(searchTools, nil, nil, nil, p.content))
}

func (p *PlaylistPage) submitSearchQuery(title string) {

}
