package ui

import (
	stdcontext "context"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"meowyplayer/scrapers"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ExplorePage struct {
	widget.BaseWidget
	userContext   *storages.UserContext
	searchEngine  scrapers.MusicSearcher
	searchResults []scrapers.Result

	cancelSearch  stdcontext.CancelFunc
	searchMutex   sync.Mutex
	debounceTimer *time.Timer

	searchEntry  *widget.Entry
	searchButton *widget.Button
	scrollList   *widget.List
}

func newExplorePage(userContext *storages.UserContext) *ExplorePage {
	p := &ExplorePage{
		userContext:  userContext,
		searchEngine: scrapers.NewPipedSearcher(),
		searchEntry:  widget.NewEntry(),
		searchButton: widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))

	p.searchEntry.OnChanged = func(t string) {
		if p.debounceTimer != nil {
			p.debounceTimer.Stop()
		}
		p.debounceTimer = time.AfterFunc(500*time.Millisecond, func() { p.submitSearchQuery(t) })
	}

	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() {
		if p.debounceTimer != nil {
			p.debounceTimer.Stop()
		}
		go p.submitSearchQuery(p.searchEntry.Text)
	}

	p.scrollList = widget.NewList(
		func() int { return len(p.searchResults) },
		func() fyne.CanvasObject { return mwidget.NewThumbnailCard(p.openInBrowser, p.showAddToPlaylistsDialog) },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*mwidget.ThumbnailCard).Set(p.searchResults[id])
		},
	)

	p.ExtendBaseWidget(p)
	return p
}

func (p *ExplorePage) openInBrowser(result scrapers.Result) {
	switch result.Platform {
	case storages.YouTubeSource:
		url, err := url.Parse("https://www.youtube.com/watch?v=" + result.ID)
		if err != nil {
			slog.Error("failed to parse url", "error", err, "id", result.ID)
			return
		}
		if err = fyne.CurrentApp().OpenURL(url); err != nil {
			slog.Error("failed to open url in browser", "error", err, "ur", url)
			return
		}
	default:
		slog.Error("unsupported platform", "platform", result.Platform)
		return
	}
}

func (p *ExplorePage) showAddToPlaylistsDialog(res scrapers.Result) {
	playlists, err := p.userContext.GetPlaylistsFromUser()
	if err != nil {
		slog.Error("failed to list playlists", "error", err)
		return
	}

	opts := make([]string, len(playlists))
	for i, pl := range playlists {
		opts[i] = pl.Title
	}

	sel := widget.NewSelect(opts, nil)
	sel.PlaceHolder = lang.L("Select a playlist")

	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowCustomConfirm(lang.L("Add to playlist"), lang.L("Add"), lang.L("Cancel"), sel, func(confirm bool) {
		i := sel.SelectedIndex()
		if !confirm || i == -1 {
			return
		}

		m := storages.Music{
			MusicId:       res.ID,
			Source:        res.Platform,
			Title:         res.Title,
			LengthSeconds: int64(res.Length.Seconds()),
		}

		if err := p.userContext.PutMusic(m); err == nil {
			if err := p.userContext.PutMusicInPlaylist(playlists[i].PlaylistId, m.MusicId, m.Source); err != nil {
				slog.Error("failed to put music in playlist", "error", err)
			}
		} else {
			slog.Error("failed to put music", "error", err)
		}
	}, win)
}

func (p *ExplorePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		mcontainer.NewHSplit(0.20, layout.NewSpacer(), p.searchEntry),
		nil, nil, nil,
		mcontainer.NewCenter(0.8, 1, p.scrollList)))
}

func (p *ExplorePage) submitSearchQuery(query string) {
	if query == "" {
		return
	}

	p.searchMutex.Lock()
	if p.cancelSearch != nil {
		p.cancelSearch()
	}
	ctx, cancel := stdcontext.WithCancel(stdcontext.Background())
	p.cancelSearch = cancel
	p.searchMutex.Unlock()

	if results, err := p.searchEngine.Search(ctx, query); err == nil {
		p.searchResults = results
		fyne.Do(func() { p.scrollList.ScrollToTop(); p.Refresh() })
	} else if ctx.Err() == nil {
		slog.Error("scraper failed", "query", query, "error", err)
		fyne.Do(func() { dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0]) })
	}
}
