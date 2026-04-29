package ui

import (
	stdcontext "context"
	"fmt"
	"log/slog"
	"net/url"
	"sync"

	"meowyplayer/context"
	"meowyplayer/scrapers"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ExplorePage struct {
	widget.BaseWidget
	searchEntry    *widget.Entry
	searchButton   *widget.Button
	scrollList     *widget.List
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

	p.scrollList = widget.NewList(
		func() int {
			return len(p.searchResults)
		},
		func() fyne.CanvasObject {
			return mwidget.NewThumbnailCard(p.openInBrowserCallback, p.addToPlaylistCallback)
		},
		func(index widget.ListItemID, object fyne.CanvasObject) {
			object.(*mwidget.ThumbnailCard).Set(p.searchResults[index])
		},
	)

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *ExplorePage) openInBrowserCallback(result scrapers.Result) {
	switch result.Platform {
	case storages.YouTubeSource:
		url, err := url.Parse(fmt.Sprintf("https://www.youtube.com/watch?v=%v", result.ID))
		if err != nil {
			slog.Error("failed to parse url", "error", err, "id", result.ID)
			return
		}
		err = fyne.CurrentApp().OpenURL(url)
		if err != nil {
			slog.Error("failed to open url in browser", "error", err, "ur", url)
			return
		}
	default:
		slog.Error("unsupported platform", "platform", result.Platform)
		return
	}
}

func (p *ExplorePage) addToPlaylistCallback(result scrapers.Result) {
	playlists, err := p.userContext.Storage().GetAllPlaylists()
	if err != nil {
		slog.Error("failed to list the playlists", "error", err)
		return
	}

	options := make([]string, 0, len(playlists))
	for _, playlist := range playlists {
		options = append(options, playlist.Title)
	}
	selects := widget.NewSelect(options, nil)
	selects.PlaceHolder = lang.L("Select a playlist")

	fyne.Do(func() {
		dialog.ShowCustomConfirm(lang.L("Add to playlist"), lang.L("Add"), lang.L("Cancel"), selects,
			func(confirm bool) {
				if i := selects.SelectedIndex(); i != -1 && confirm {
					go func() {
						music := storages.Music{
							MusicId:       result.ID,
							Source:        result.Platform,
							Title:         result.Title,
							LengthSeconds: int64(result.Length.Seconds()),
						}

						// Download the music if not in the local storage.
						musicFile, err := p.userContext.Storage().GetMusicFile(music)
						if err != nil {
							// Download the missing music files.
							content, err := p.downloadEngine.Download(stdcontext.Background(), result)
							if err != nil {
								slog.Error("failed to download music", "error", err)
								return
							}
							defer content.Close()
							err = p.userContext.Storage().CreateOrUpdateMusicFile(music, content)
							if err != nil {
								slog.Error("failed to create music file", "error", err)
								return
							}
						}
						musicFile.Close()

						// Add music relation to DB.
						err = p.userContext.Storage().CreateMusic(music)
						if err != nil {
							slog.Error("failed to create the music", "error", err)
							return
						}

						err = p.userContext.AddMusicToPlaylist(playlists[i], music.MusicId, music.Source)
						if err != nil {
							slog.Error("failed to add music to the playlist", "error", err)
							return
						}
					}()
				}
			}, fyne.CurrentApp().Driver().AllWindows()[0])
	})
}

func (p *ExplorePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		mcontainer.NewCenter(0.62, 1, p.searchEntry), nil, nil, nil,
		mcontainer.NewCenter(0.8, 1, p.scrollList)))
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
	p.searchResults = results

	fyne.Do(func() {
		p.scrollList.ScrollToTop()
		p.Refresh()
	})
}
