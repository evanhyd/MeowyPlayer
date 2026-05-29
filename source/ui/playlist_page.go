package ui

import (
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mutil"
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
	userContext          *context.UserContext
	queryResults         []storages.Playlist
	displayResults       []storages.Playlist
	searchEntry          *widget.Entry
	searchButton         *widget.Button
	scrollList           *widget.GridWrap
	createPlaylistButton *widget.Button
}

func newPlaylistPage(userContext *context.UserContext) *PlaylistPage {
	p := PlaylistPage{
		userContext:          userContext,
		searchEntry:          widget.NewEntry(),
		searchButton:         widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		createPlaylistButton: widget.NewButtonWithIcon(lang.L("Create Playlist"), theme.FolderNewIcon(), nil),
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.filterResults
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.filterResults(p.searchEntry.Text) }
	p.createPlaylistButton.Importance = widget.LowImportance
	p.createPlaylistButton.OnTapped = p.showCreatePlaylistDialog

	p.scrollList = widget.NewGridWrap(
		func() int {
			return len(p.displayResults)
		},
		func() fyne.CanvasObject {
			return mwidget.NewPlaylistCard(
				func(playlist storages.Playlist) {
					p.Hide()
					p.userContext.ViewMusicPage(playlist)
				},
				p.showEditingMenu,
			)
		},
		func(index widget.GridWrapItemID, object fyne.CanvasObject) {
			object.(*mwidget.PlaylistCard).Set(p.displayResults[index])
		},
	)

	p.userContext.AddListener(context.OnSetStorageEvent, func(any) {
		p.fetchPlaylists()
		p.Show()
	})
	p.userContext.AddListener(context.OnViewPlaylistPageEvent, func(any) {
		p.Show()
	})
	p.userContext.AddListener(context.OnPutPlaylistEvent, func(any) {
		p.fetchPlaylists()
	})
	p.userContext.AddListener(context.OnDeletePlaylistEvent, func(any) {
		p.fetchPlaylists()
	})

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *PlaylistPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		mcontainer.NewCenter(0.62, 1, container.NewBorder(nil, nil, nil, p.createPlaylistButton, p.searchEntry)),
		nil,
		nil,
		nil,
		p.scrollList))
}

func (p *PlaylistPage) showEditingMenu(playlist storages.Playlist, event *fyne.PointEvent) {
	editMenu := fyne.NewMenuItemWithIcon(lang.L("Edit"), theme.DocumentCreateIcon(), func() {
		p.showEditPlaylistDialog(playlist)
	})
	deleteMenu := fyne.NewMenuItemWithIcon(lang.L("Delete"), theme.DeleteIcon(), func() {
		p.showDeletePlaylistDialog(playlist)
	})
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), fyne.CurrentApp().Driver().AllWindows()[0].Canvas(), event.AbsolutePosition)
}

func (p *PlaylistPage) showEditPlaylistDialog(playlist storages.Playlist) {
	editor := newPlaylistEditorWithState(playlist.Title, fyne.NewStaticResource(mutil.PlaylistIdToString(playlist.PlaylistId), playlist.CoverBlob))
	dialog.ShowCustomConfirm(lang.L("Edit Album"), lang.L("save"), lang.L("cancel"), editor,
		func(confirm bool) {
			if confirm {
				playlist.Title, playlist.CoverBlob = editor.state()
				if _, err := p.userContext.PutPlaylist(playlist); err != nil {
					slog.Error("failed to update the playlist", "error", err)
					return
				}
			}
		},
		fyne.CurrentApp().Driver().AllWindows()[0],
	)
}

func (p *PlaylistPage) showDeletePlaylistDialog(playlist storages.Playlist) {
	dialog.ShowCustomConfirm(lang.L("Delete Playlist Confirmation"), lang.L("delete"), lang.L("cancel"),
		widget.NewLabel(lang.L("Do you want to delete the playlist: ")+playlist.Title),
		func(confirm bool) {
			if confirm {
				if err := p.userContext.DeletePlaylist(playlist.PlaylistId); err != nil {
					slog.Error("failed to delete the playlist", "error", err)
					return
				}
			}
		},
		fyne.CurrentApp().Driver().AllWindows()[0],
	)
}

func (p *PlaylistPage) showCreatePlaylistDialog() {
	editor := newPlaylistEditor()
	dialog.ShowCustomConfirm(lang.L("Create Playlist"), lang.L("Create"), lang.L("Cancel"), editor,
		func(confirm bool) {
			if confirm {
				title, coverBlob := editor.state()
				_, err := p.userContext.PutPlaylist(storages.Playlist{Title: title, CoverBlob: coverBlob})
				if err != nil {
					slog.Error("failed to create playlist", "error", err)
				}
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0],
	)
}

func (p *PlaylistPage) fetchPlaylists() {
	var err error
	p.queryResults, err = p.userContext.GetPlaylistsFromUser()
	if err != nil {
		slog.Error("failed to query playlists", "error", err)
		return
	}
	p.filterResults(p.searchEntry.Text)
}

func (p *PlaylistPage) filterResults(title string) {
	// Filter by title.
	title = strings.ToLower(title)
	p.displayResults = p.displayResults[:0]
	for i := range p.queryResults {
		if strings.Contains(strings.ToLower(p.queryResults[i].Title), title) {
			p.displayResults = append(p.displayResults, p.queryResults[i])
		}
	}

	fyne.Do(func() {
		p.scrollList.ScrollToTop()
		p.scrollList.Refresh()
	})
}
