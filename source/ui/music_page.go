package ui

import (
	"fmt"
	"image/color"
	"log/slog"
	"meowyplayer/mcontext"
	"meowyplayer/scrapers"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mutil"
	"meowyplayer/ui/internal/mwidget"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicPage struct {
	widget.BaseWidget
	userContext    *mcontext.UserContext
	playlist       storages.Playlist
	queryResults   []storages.Music
	displayResults []storages.Music
	downloader     scrapers.MusicDownloader

	background          *canvas.Image
	fadeOverlay         *canvas.LinearGradient
	searchEntry         *widget.Entry
	searchButton        *widget.Button
	backButton          *widget.Button
	playlistCover       *canvas.Image
	playlistTitle       *widget.Label
	playlistDescription *widget.Label
	scrollList          *widget.List
}

func newMusicPage(userContext *mcontext.UserContext) *MusicPage {
	p := MusicPage{
		userContext:         userContext,
		background:          canvas.NewImageFromImage(nil),
		fadeOverlay:         canvas.NewVerticalGradient(color.Transparent, theme.Color(theme.ColorNameBackground)),
		searchEntry:         widget.NewEntry(),
		searchButton:        widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		backButton:          widget.NewButtonWithIcon(lang.L("Back"), theme.NavigateBackIcon(), nil),
		playlistCover:       canvas.NewImageFromResource(nil),
		playlistTitle:       widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		playlistDescription: widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
	}

	p.background.ScaleMode = canvas.ImageScaleFastest

	p.searchEntry.ActionItem = p.searchButton
	p.searchEntry.SetPlaceHolder(lang.L("Search songs, videos, or artists"))
	p.searchEntry.OnChanged = p.filterResults

	p.searchButton.Importance = widget.LowImportance
	p.searchButton.OnTapped = func() { p.filterResults(p.searchEntry.Text) }

	p.backButton.Importance = widget.LowImportance
	p.backButton.OnTapped = p.userContext.ViewPlaylistPage

	p.playlistCover.CornerRadius = 8.0
	p.playlistCover.SetMinSize(mwidget.PlaylistCardSize)
	p.playlistCover.FillMode = canvas.ImageFillContain
	p.playlistCover.ScaleMode = canvas.ImageScaleFastest

	p.playlistTitle.Truncation = fyne.TextTruncateEllipsis
	p.playlistDescription.Truncation = fyne.TextTruncateEllipsis

	p.scrollList = widget.NewList(
		func() int {
			return len(p.displayResults)
		},
		func() fyne.CanvasObject {
			return mwidget.NewMusicCard(
				func(music storages.Music) {
					p.userContext.PlayMusic(p.playlist, music)
				},
				p.showEditingMenu,
			)
		},
		func(index widget.GridWrapItemID, object fyne.CanvasObject) {
			object.(*mwidget.MusicCard).Set(p.displayResults[index])
		},
	)
	p.scrollList.HideSeparators = true

	p.userContext.AddListener(mcontext.OnSetStorageEvent, func(any) {
		p.Hide()
	})

	p.userContext.AddListener(mcontext.OnViewMusicPageEvent, func(data any) {
		p.fetchMusic(data.(mcontext.OnViewMusicPageEventData).Playlist)
		p.Show()
	})

	p.userContext.AddListener(mcontext.OnViewPlaylistPageEvent, func(any) {
		p.Hide()
	})

	p.userContext.AddListener(mcontext.OnPutMusicInPlaylistEvent, func(data any) {
		if data.(mcontext.OnPutMusicInPlaylistEventData).PlaylistId == p.playlist.PlaylistId {
			p.fetchMusic(p.playlist)
		}
	})

	p.userContext.AddListener(mcontext.OnDeleteMusicFromPlaylistEvent, func(data any) {
		p.fetchMusic(p.playlist)
	})

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *MusicPage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		mcontainer.NewVSplit(0.5, container.NewStack(p.background, p.fadeOverlay), layout.NewSpacer()),
		container.NewBorder(
			mcontainer.NewHSplit(0.20,
				layout.NewSpacer(),
				container.NewBorder(nil, nil, nil, p.backButton, p.searchEntry),
			),
			nil,
			nil,
			nil,
			mcontainer.NewHSplit(0.35,
				mcontainer.NewCenter(0.85, 0.85, mcontainer.NewVSplit(0.5,
					p.playlistCover,
					container.NewVBox(p.playlistTitle, p.playlistDescription),
				)),
				p.scrollList,
			),
		),
	))
}

func (p *MusicPage) showEditingMenu(music storages.Music, event *fyne.PointEvent) {
	editMenu := fyne.NewMenuItemWithIcon(lang.L("Details"), theme.DocumentCreateIcon(), func() {
		p.showDetailDialog(music)
	})
	deleteMenu := fyne.NewMenuItemWithIcon(lang.L("Delete"), theme.DeleteIcon(), func() {
		p.showDeleteMusicDialog(music)
	})
	widget.ShowPopUpMenuAtPosition(fyne.NewMenu("", editMenu, deleteMenu), fyne.CurrentApp().Driver().AllWindows()[0].Canvas(), event.AbsolutePosition)
}

func (p *MusicPage) showDetailDialog(music storages.Music) {
	fyne.CurrentApp().Clipboard().SetContent(music.MusicId)
	fyne.CurrentApp().SendNotification(fyne.NewNotification("Success", lang.L("Successfully copied music ID.")))

	dialog.ShowInformation(lang.L("Music Detail"),
		fmt.Sprintf("%v: %v\nID: %v\n%v: %v\n%v: %v",
			lang.L("Title"), music.Title,
			music.MusicId,
			lang.L("Platform"), func() string {
				switch music.Source {
				case storages.UnknownSource:
					return "Unknown"
				case storages.YouTubeSource:
					return "YouTube"
				case storages.SpotifySource:
					return "Spotify"
				default:
					return "Error"
				}
			}(),
			lang.L("Duration"), mutil.SecondsToTime(music.LengthSeconds),
		), fyne.CurrentApp().Driver().AllWindows()[0])
}

func (p *MusicPage) showDeleteMusicDialog(music storages.Music) {
	dialog.ShowCustomConfirm(lang.L("Delete Music Confirmation"), lang.L("Delete"), lang.L("Cancel"),
		widget.NewLabel(lang.L("Do you want to delete ")+music.Title+" from the playlist"),
		func(confirm bool) {
			if confirm {
				if err := p.userContext.DeleteMusicFromPlaylist(p.playlist.PlaylistId, music.MusicId, music.Source); err != nil {
					slog.Error("failed to delete the music from the playlist", "error", err)
					return
				}
			}
		},
		fyne.CurrentApp().Driver().AllWindows()[0],
	)
}

func (p *MusicPage) fetchMusic(playlist storages.Playlist) {
	p.playlist = playlist

	var err error
	p.background.Image, err = mutil.ScaleImageFromBytes(playlist.CoverBlob, fyne.NewSize(8, 8))
	if err != nil {
		slog.Error("failed to blur images", "error", err)
		return
	}
	p.background.Translucency = 0.8
	p.background.Refresh()
	p.playlistCover.Resource = fyne.NewStaticResource(mutil.PlaylistIdToString(playlist.PlaylistId), playlist.CoverBlob)
	p.playlistCover.Refresh()
	p.playlistTitle.SetText(playlist.Title)

	p.queryResults, err = p.userContext.GetMusicFromPlaylist(p.playlist.PlaylistId)
	if err != nil {
		slog.Error("failed to query music in playlist", "error", err)
		return
	}

	var totalSeconds int64
	for i := range p.queryResults {
		totalSeconds += p.queryResults[i].LengthSeconds
	}

	p.playlistDescription.SetText(
		fmt.Sprintf("%v %v - %v %v\n%v: %v",
			len(p.queryResults), lang.L("songs"), mutil.SecondsToTime(totalSeconds), lang.L("minutes"),
			lang.L("Modified"), time.Unix(0, playlist.ModifiedDate).Format(time.DateTime),
		))

	p.filterResults(p.searchEntry.Text)
}

func (p *MusicPage) filterResults(title string) {
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
