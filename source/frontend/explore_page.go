package frontend

import (
	"log/slog"
	"meowyplayer/frontend/internal/layouts"
	"meowyplayer/frontend/internal/widgets"
	"meowyplayer/scrapers"

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
	searchButton *widget.Button
	searchEntry  *widget.Entry
	content      *widget.List

	scraper       scrapers.MusicScraper
	searchResults []scrapers.Result
}

func newExplorePage() *ExplorePage {
	p := ExplorePage{
		searchButton: widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		searchEntry:  widget.NewEntry(),

		scraper: scrapers.NewYouTubeScraper(),
	}

	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.submitSearchQuery(p.searchEntry.Text) }

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnSubmitted = p.submitSearchQuery

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
	searchRow := container.New(layouts.NewHSegmentLayout(2, 7, 2), layout.NewSpacer(), p.searchEntry, layout.NewSpacer())
	return widget.NewSimpleRenderer(container.NewBorder(
		searchRow, nil, nil, nil, container.New(layouts.NewHSegmentLayout(1, 8, 1), layout.NewSpacer(), p.content, layout.NewSpacer()),
	))
}

func (p *ExplorePage) submitSearchQuery(query string) {
	results, err := p.scraper.Search(query)
	if err != nil {
		slog.Error("scraper failed to search the query", "query", query, "error", err)
		dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
	}

	p.searchResults = results
	p.content.Refresh()
}
