package ui

import (
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mwidget"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlaylistPage struct {
	widget.BaseWidget
	searchEntry          *widget.Entry
	searchButton         *widget.Button
	scrollList           *widget.GridWrap
	createPlaylistButton *widget.Button
	queryResults         []storages.Playlist
	displayResults       []storages.Playlist
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
	p.searchEntry.OnChanged = p.updateDisplayResults
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.updateDisplayResults(p.searchEntry.Text) }
	p.createPlaylistButton.Importance = widget.LowImportance
	p.createPlaylistButton.OnTapped = p.createPlaylist

	p.scrollList = widget.NewGridWrap(
		func() int {
			return len(p.displayResults)
		},
		func() fyne.CanvasObject {
			return mwidget.NewPlaylistCard(func(playlist storages.Playlist) {
				p.Hide()
				p.userContext.ViewPlaylist(playlist)
			})
		},
		func(index widget.GridWrapItemID, object fyne.CanvasObject) {
			object.(*mwidget.PlaylistCard).Set(p.displayResults[index])
		},
	)

	p.userContext.AddListener(context.OnSetStorageEvent, func(context.EventType, any) {
		p.fetchPlaylists()
		p.Show()
	})
	p.userContext.AddListener(context.OnCreatePlaylistEvent, func(context.EventType, any) {
		p.fetchPlaylists()
	})
	p.userContext.AddListener(context.OnReturnBackFromPlaylistEvent, func(context.EventType, any) {
		p.Show()
	})

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *PlaylistPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		container.NewStack(mcontainer.NewCenter(0.62, 1, p.searchEntry), container.NewBorder(nil, nil, nil, p.createPlaylistButton)),
		nil, nil, nil, p.scrollList))
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

func (p *PlaylistPage) fetchPlaylists() {
	var err error
	p.queryResults, err = p.userContext.Storage().GetAllPlaylists()
	if err != nil {
		slog.Error("failed to query playlists", "error", err)
		return
	}
	p.updateDisplayResults(p.searchEntry.Text)
}

func (p *PlaylistPage) updateDisplayResults(title string) {
	// Filter by title.
	title = strings.ToLower(title)
	p.displayResults = p.displayResults[:0]
	for _, playlist := range p.queryResults {
		if strings.Contains(strings.ToLower(playlist.Title), title) {
			p.displayResults = append(p.displayResults, playlist)
		}
	}

	fyne.Do(func() {
		p.scrollList.ScrollToTop()
		p.scrollList.Refresh()
	})
}
