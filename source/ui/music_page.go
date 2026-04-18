package ui

import (
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/layouts"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicPage struct {
	widget.BaseWidget
	searchEntry  *widget.Entry
	searchButton *widget.Button
	content      *widget.GridWrap
	backButton   *widget.Button

	playlist       storages.Playlist
	queryResults   []storages.Music
	displayResults []storages.Music
	userContext    *context.UserContext
}

func newMusicPage(userContext *context.UserContext) *MusicPage {
	p := MusicPage{
		searchEntry:  widget.NewEntry(),
		searchButton: widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		backButton:   widget.NewButtonWithIcon(lang.L("Back"), theme.NavigateBackIcon(), nil),
		userContext:  userContext,
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.updateDisplayResults
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.updateDisplayResults(p.searchEntry.Text) }
	p.backButton.Importance = widget.LowImportance
	p.backButton.OnTapped = p.userContext.ReturnBackFromPlaylist

	// TODO: implement
	// p.content = widget.NewGridWrap(
	// 	func() int {
	// 		return len(p.displayResults)
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return widgets.NewPlaylistCard(func(playlistId int64) {})
	// 	},
	// 	func(index widget.GridWrapItemID, object fyne.CanvasObject) {
	// 		object.(*widgets.PlaylistCard).Set(p.displayResults[index])
	// 	},
	// )

	p.userContext.AddListener(context.OnViewPlaylistEvent, func(context.EventType, any) {
		p.Show()
	})

	p.userContext.AddListener(context.OnSetStorageEvent, func(context.EventType, any) {
		p.Hide()
	})

	p.userContext.AddListener(context.OnReturnBackFromPlaylistEvent, func(context.EventType, any) {
		p.Hide()
	})

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *MusicPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		container.New(layouts.NewCenterLayout(0.62, 1), container.NewBorder(nil, nil, nil, p.backButton, p.searchEntry)), nil, nil, nil, widget.NewButtonWithIcon("button", theme.BrokenImageIcon(), nil)))
}

func (p *MusicPage) fetchMusic(context.EventType, any) {
	var err error
	p.queryResults, err = p.userContext.Storage().GetAllMusicFromPlaylist(p.playlist.PlaylistId)
	if err != nil {
		slog.Error("failed to query music in playlist", "error", err)
		return
	}
	p.updateDisplayResults(p.searchEntry.Text)
}

func (p *MusicPage) updateDisplayResults(title string) {
	// Filter by title.
	title = strings.ToLower(title)
	p.displayResults = p.displayResults[:0]
	for _, music := range p.queryResults {
		if strings.Contains(strings.ToLower(music.Title), title) {
			p.displayResults = append(p.displayResults, music)
		}
	}

	fyne.Do(func() {
		p.content.ScrollToTop()
		p.content.Refresh()
	})
}
