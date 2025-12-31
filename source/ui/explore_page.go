package ui

import (
	"log/slog"
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
	searchEntry   *widget.Entry
	searchButton  *widget.Button
	content       *widget.List
	searchResults []scrapers.Result

	searchEngine   scrapers.MusicSearcher
	downloadEngine scrapers.MusicDownloader
}

func newExplorePage() *ExplorePage {
	p := ExplorePage{
		searchEntry:    widget.NewEntry(),
		searchButton:   widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		searchEngine:   scrapers.NewClipzagSearche(),
		downloadEngine: scrapers.NewCnvmp3Downloader(),
	}

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.submitSearchQuery
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.submitSearchQuery(p.searchEntry.Text) }

	p.content = widget.NewList(
		func() int {
			return len(p.searchResults)
		},
		func() fyne.CanvasObject {
			return widgets.NewThumbnailCard()
		},
		func(index widget.ListItemID, object fyne.CanvasObject) {
			object.(*widgets.ThumbnailCard).SetResult(p.searchResults[index])
		},
	)

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *ExplorePage) CreateRenderer() fyne.WidgetRenderer {
	searchTools := container.New(layouts.NewHSegmentLayout(2, 7, 2), layout.NewSpacer(), p.searchEntry, layout.NewSpacer())
	return widget.NewSimpleRenderer(container.NewBorder(searchTools, nil, nil, nil, p.content))
}

func (p *ExplorePage) submitSearchQuery(query string) {
	go func() {
		results, err := p.searchEngine.Search(query)
		if err != nil {
			slog.Error("scraper failed to search the query", "query", query, "error", err)
			dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
		}

		p.searchResults = results
		p.content.Refresh()
		p.content.ScrollToTop()
	}()
}
