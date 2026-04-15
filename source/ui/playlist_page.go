package ui

import (
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/layouts"
	"meowyplayer/ui/internal/widgets"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	userContext   *context.UserContext
}

func newPlaylistPage(context *context.UserContext) *PlaylistPage {
	p := PlaylistPage{
		searchEntry:  widget.NewEntry(),
		searchButton: widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		userContext:  context,
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.submitSearchQuery
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.submitSearchQuery(p.searchEntry.Text) }

	p.content = widget.NewGridWrap(
		func() int {
			return len(p.searchResults)
		},
		func() fyne.CanvasObject {
			return widgets.NewPlaylistCard(func(playlistId int64) {})
		},
		func(index widget.GridWrapItemID, object fyne.CanvasObject) {
			object.(*widgets.PlaylistCard).Set(p.searchResults[index])
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
	fyne.Do(func() {
		playlists, err := p.userContext.Storage().ListAllPlaylists()
		if err != nil {
			slog.Error("failed to query playlists", "title", title, "error", err)
			dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
			return
		}

		// Filter by title.
		title = strings.ToLower(title)
		playlists = slices.DeleteFunc(playlists, func(p storages.Playlist) bool {
			return !strings.Contains(strings.ToLower(p.Title), title)
		})

		p.searchResults = playlists
		p.content.Refresh()
		p.content.ScrollToTop()
	})
}
