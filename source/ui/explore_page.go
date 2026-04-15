package ui

import (
	stdcontext "context"
	"log/slog"
	"sync"

	"meowyplayer/context"
	"meowyplayer/scrapers"
	"meowyplayer/ui/internal/layouts"
	"meowyplayer/ui/internal/widgets"

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
	searchEntry    *widget.Entry
	searchButton   *widget.Button
	content        *widget.List
	searchResults  []scrapers.Result
	searchEngine   scrapers.MusicSearcher
	downloadEngine scrapers.MusicDownloader

	userContext  *context.UserContext
	cancelSearch stdcontext.CancelFunc
	searchMutex  sync.Mutex
}

func newExplorePage(userContext *context.UserContext) *ExplorePage {
	p := ExplorePage{
		searchEntry:    widget.NewEntry(),
		searchButton:   widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		searchEngine:   scrapers.NewInvidiousSearcher(),
		downloadEngine: scrapers.NewCnvmp3Downloader(),
		userContext:    userContext,
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = func(title string) { go p.submitSearchQuery(title) }
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { go p.submitSearchQuery(p.searchEntry.Text) }

	p.content = widget.NewList(
		func() int {
			return len(p.searchResults)
		},
		func() fyne.CanvasObject {
			return widgets.NewThumbnailCard()
		},
		func(index widget.ListItemID, object fyne.CanvasObject) {
			object.(*widgets.ThumbnailCard).Set(p.searchResults[index])
		},
	)

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *ExplorePage) CreateRenderer() fyne.WidgetRenderer {
	searchTools := container.New(layouts.NewHSegmentLayout(2, 7, 2), layout.NewSpacer(), p.searchEntry, layout.NewSpacer())
	return widget.NewSimpleRenderer(container.NewBorder(
		searchTools, nil, nil, nil,
		container.New(layouts.NewHSegmentLayout(1, 7, 1), layout.NewSpacer(), p.content, layout.NewSpacer())))
}

func (p *ExplorePage) submitSearchQuery(query string) {
	// Set up search cancelling.
	p.searchMutex.Lock()
	if p.cancelSearch != nil {
		p.cancelSearch()
	}
	ctx, cancel := stdcontext.WithCancel(stdcontext.Background())
	p.cancelSearch = cancel
	p.searchMutex.Unlock()

	// Handle empty queries.
	if query == "" {
		return
	}

	results, err := p.searchEngine.Search(ctx, query)
	if ctx.Err() != nil {
		return
	}

	if err != nil {
		slog.Error("scraper failed to search the query", "query", query, "error", err)
		fyne.Do(func() {
			dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
		})
		return
	}

	fyne.Do(func() {
		p.searchResults = results
		p.content.Refresh()
		p.content.ScrollToTop()
	})
}
