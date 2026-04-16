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
	searchEntry          *widget.Entry
	searchButton         *widget.Button
	content              *widget.GridWrap
	createPlaylistButton *widget.Button
	searchResults        []storages.Playlist
	userContext          *context.UserContext
}

func newPlaylistPage(userContext *context.UserContext) *PlaylistPage {
	p := PlaylistPage{
		searchEntry:          widget.NewEntry(),
		searchButton:         widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		createPlaylistButton: widget.NewButtonWithIcon(lang.L("Create Playlist"), theme.FolderNewIcon(), nil),
		userContext:          userContext,
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = func(s string) { go p.submitSearchQuery(s) }
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { go p.submitSearchQuery(p.searchEntry.Text) }
	p.createPlaylistButton.Importance = widget.LowImportance
	p.createPlaylistButton.OnTapped = p.createPlaylist

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

	p.userContext.AddListener(context.OnSetStorageEvent, p.refreshPlaylist)
	p.userContext.AddListener(context.OnCreatePlaylistEvent, p.refreshPlaylist)

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *PlaylistPage) CreateRenderer() fyne.WidgetRenderer {
	searchTools := container.New(layouts.NewHSegmentLayout(2, 7, 2),
		layout.NewSpacer(),
		container.NewBorder(nil, nil, nil, p.createPlaylistButton, p.searchEntry),
		layout.NewSpacer())
	return widget.NewSimpleRenderer(container.NewBorder(searchTools, nil, nil, nil, p.content))
}

func (p *PlaylistPage) createPlaylist() {
	editor := newPlaylistEditor()
	dialog.ShowCustomConfirm(lang.L("Create Playlist"), lang.L("Create"), lang.L("Cancel"), editor,
		func(confirm bool) {
			if confirm {
				err := p.userContext.CreatePlaylist(editor.state())
				if err != nil {
					slog.Error("failed to create playlist", "error", err)
				}
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0],
	)
}

func (p *PlaylistPage) refreshPlaylist(context.EventType, any) {
	p.submitSearchQuery(p.searchEntry.Text)
}

func (p *PlaylistPage) submitSearchQuery(title string) {
	playlists, err := p.userContext.Storage().ListAllPlaylists()
	if err != nil {
		slog.Error("failed to query playlists", "title", title, "error", err)
		return
	}

	// Filter by title.
	title = strings.ToLower(title)
	playlists = slices.DeleteFunc(playlists, func(p storages.Playlist) bool {
		return !strings.Contains(strings.ToLower(p.Title), title)
	})

	fyne.Do(func() {
		p.searchResults = playlists
		p.content.Refresh()
		p.content.ScrollToTop()
	})
}
