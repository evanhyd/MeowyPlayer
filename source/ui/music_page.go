package ui

import (
	"image/color"
	"log/slog"
	"meowyplayer/context"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mwidget"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicPage struct {
	widget.BaseWidget

	background  *canvas.Image
	fadeOverlay *canvas.LinearGradient

	searchEntry         *widget.Entry
	searchButton        *widget.Button
	backButton          *widget.Button
	playlistCover       *canvas.Image
	playlistTitle       *widget.Label
	playlistDescription *widget.Label
	scrollList          *widget.List

	playlist       storages.Playlist
	queryResults   []storages.Music
	displayResults []storages.Music
	userContext    *context.UserContext
}

func newMusicPage(userContext *context.UserContext) *MusicPage {
	p := MusicPage{
		background:          canvas.NewImageFromImage(nil),
		fadeOverlay:         canvas.NewVerticalGradient(color.Transparent, theme.Color(theme.ColorNameBackground)),
		searchEntry:         widget.NewEntry(),
		searchButton:        widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		backButton:          widget.NewButtonWithIcon(lang.L("Back"), theme.NavigateBackIcon(), nil),
		playlistCover:       canvas.NewImageFromResource(nil),
		playlistTitle:       widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		playlistDescription: widget.NewLabel(""),
		userContext:         userContext,
	}

	p.background.ScaleMode = canvas.ImageScaleFastest
	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.updateDisplayResults
	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.updateDisplayResults(p.searchEntry.Text) }
	p.backButton.Importance = widget.LowImportance
	p.backButton.OnTapped = p.userContext.ReturnBackFromPlaylist
	p.playlistCover.CornerRadius = 8.0

	p.scrollList = widget.NewList(
		func() int {
			return len(p.displayResults)
		},
		func() fyne.CanvasObject {
			return mwidget.NewMusicCard(func(music storages.Music) {
				// TODO: request to play
			})
		},
		func(index widget.GridWrapItemID, object fyne.CanvasObject) {
			object.(*mwidget.MusicCard).Set(p.displayResults[index])
		},
	)

	p.userContext.AddListener(context.OnViewPlaylistEvent, func(_ context.EventType, data any) {
		p.fetchMusic(data.(context.OnViewPlaylistEventData).Playlist)
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
	return widget.NewSimpleRenderer(container.NewStack(
		mcontainer.NewVSplit(0.5, container.NewStack(p.background, p.fadeOverlay), layout.NewSpacer()),
		container.NewBorder(
			mcontainer.NewCenter(0.62, 1, container.NewBorder(nil, nil, nil, p.backButton, p.searchEntry)),
			nil,
			nil,
			nil,
			mcontainer.NewHSplit(0.35,
				mcontainer.NewCenter(0.8, 0.8, mcontainer.NewVSplit(0.5,
					p.playlistCover,
					container.NewVBox(p.playlistTitle, p.playlistDescription),
				)),
				p.scrollList,
			),
		),
	))
}

func (p *MusicPage) fetchMusic(playlist storages.Playlist) {
	p.playlist = playlist

	var err error
	p.background.Image, err = mwidget.ScaleImageFromBytes(playlist.CoverBlob, fyne.NewSize(8, 8))
	if err != nil {
		slog.Error("failed to blur images", "error", err)
		return
	}
	p.background.Translucency = 0.8
	p.background.Refresh()
	p.playlistCover.Resource = fyne.NewStaticResource(strconv.FormatInt(playlist.PlaylistId, 16), p.playlist.CoverBlob)
	p.playlistCover.Refresh()
	p.playlistTitle.SetText(playlist.Title)
	p.playlistDescription.SetText("Description")

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
		p.scrollList.ScrollToTop()
		p.scrollList.Refresh()
	})
}
