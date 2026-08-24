package ui

import (
	stdcontext "context"
	"fmt"
	"log/slog"
	"meowyplayer/mcontext"
	"meowyplayer/players"
	"meowyplayer/scrapers"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mutil"
	"meowyplayer/ui/internal/mwidget"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
)

const (
	Pause  = 0
	Resume = 1
)

type MusicController struct {
	widget.BaseWidget
	userContext *mcontext.UserContext
	musicPlayer players.MusicPlayer
	playlist    storages.Playlist

	playlistCover  *canvas.Image
	title          *widget.RichText
	progressSlider *mwidget.ProgressSlider
	durationLabel  *widget.Label
	modeDropDown   *mwidget.DropDown
	skipPrevButton *widget.Button
	playButton     *widget.Button
	skipNextButton *widget.Button
	volumeSlider   *mwidget.VolumeSlider
}

func newMusicController(userContext *mcontext.UserContext) *MusicController {
	var c MusicController
	c = MusicController{
		userContext:   userContext,
		musicPlayer:   players.MakeBeepPlayer(userContext, c.downloadMusic, c.onPlayMusic),
		playlistCover: canvas.NewImageFromResource(resourceIconPng),
		title: widget.NewRichText(&widget.TextSegment{
			Style: widget.RichTextStyle{SizeName: theme.SizeNameSubHeadingText, TextStyle: fyne.TextStyle{Bold: true}},
		}),
		durationLabel:  widget.NewLabel("00:00"),
		modeDropDown:   mwidget.NewDropDown(),
		skipPrevButton: widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil),
		playButton:     widget.NewButtonWithIcon("", theme.MediaRecordIcon(), nil),
		skipNextButton: widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil),
		volumeSlider:   mwidget.NewVolumeSlider(),
	}

	c.playlistCover.SetMinSize(mwidget.PlaylistCardSize)
	c.playlistCover.CornerRadius = 8.0

	c.title.Truncation = fyne.TextTruncateEllipsis

	c.progressSlider = mwidget.NewProgressSlider(c.musicPlayer.SetProgress)

	c.modeDropDown.Add(fyne.NewMenuItemWithIcon(lang.L("Sequential"), theme.MailForwardIcon(), func() { c.musicPlayer.SetQueueMode(players.SequentialQueueMode) }))
	c.modeDropDown.Add(fyne.NewMenuItemWithIcon(lang.L("Random"), resourceRandomSvg, func() { c.musicPlayer.SetQueueMode(players.RandomQueueMode) }))
	c.modeDropDown.Select(0)

	c.skipPrevButton.Importance = widget.LowImportance
	c.skipPrevButton.OnTapped = c.musicPlayer.Previous

	c.playButton.Importance = widget.LowImportance
	c.playButton.OnTapped = func() {
		if c.musicPlayer.IsPlaying() {
			c.musicPlayer.Pause()
		} else {
			c.musicPlayer.Resume()
		}
	}

	c.skipNextButton.Importance = widget.LowImportance
	c.skipNextButton.OnTapped = c.musicPlayer.Next

	c.volumeSlider.OnChanged = c.musicPlayer.SetVolume
	c.volumeSlider.SetVolume(0.7)

	c.userContext.AddListener(mcontext.OnPlayMusicEvent, func(data any) {
		c.fetchMusic(data.(mcontext.OnPlayMusicEventData).Playlist, data.(mcontext.OnPlayMusicEventData).Music)
	})

	// UI update thread
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			progress := c.musicPlayer.GetProgress()
			music := c.musicPlayer.GetMusic()
			playedDuration := int64(float64(music.LengthSeconds) * progress)
			fyne.DoAndWait(func() {
				c.durationLabel.SetText(fmt.Sprintf("%s / %s", mutil.SecondsToTime(playedDuration), mutil.SecondsToTime(music.LengthSeconds)))
				c.progressSlider.SetValue(progress)
			})
		}
	}()

	c.ExtendBaseWidget(&c)
	return &c
}

func (c *MusicController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil,
		nil,
		c.playlistCover,
		nil,
		container.NewBorder(
			c.title,
			mcontainer.NewHSplit(0.8, container.NewCenter(container.NewHBox(c.modeDropDown, c.skipPrevButton, c.playButton, c.skipNextButton)), c.volumeSlider),
			nil,
			nil,
			container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider),
		),
	))
}

func (c *MusicController) extractCover(music storages.Music) []byte {
	file, err := c.userContext.GetMusicFile(music)
	if err != nil {
		return nil
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		return nil
	}

	pic := m.Picture()
	if pic == nil {
		return nil
	}
	return pic.Data
}

func (c *MusicController) onPlayMusic(music storages.Music) {
	// Update title and cover.
	c.title.Segments[0].(*widget.TextSegment).Text = music.Title
	c.title.Refresh()

	if cover := c.extractCover(music); cover != nil {
		c.playlistCover.Resource = fyne.NewStaticResource(music.MusicId, cover)
	} else {
		c.playlistCover.Resource = fyne.NewStaticResource(music.MusicId, c.playlist.CoverBlob)
	}
	c.playlistCover.Refresh()
}

func (c *MusicController) fetchMusic(playlist storages.Playlist, music storages.Music) {
	c.playlist = playlist
	c.musicPlayer.SetPlaylist(playlist, music)
}

func (c *MusicController) downloadMusic(music storages.Music) {
	progressBar := dialog.NewCustomWithoutButtons(lang.L("Downloading"), widget.NewProgressBarInfinite(), fyne.CurrentApp().Driver().AllWindows()[0])
	fyne.DoAndWait(progressBar.Show)

	// Download the missing music files.
	downloader := scrapers.NewCnvmp3Downloader()
	content, err := downloader.Download(stdcontext.Background(), scrapers.Result{
		Platform: music.Source,
		ID:       music.MusicId,
		Title:    music.Title,
	})
	if err != nil {
		slog.Error("failed to download music", "error", err)
		return
	}
	defer content.Close()

	err = c.userContext.PutMusicFile(music, content)
	if err != nil {
		slog.Error("failed to put music file", "error", err)
		return
	}

	fyne.DoAndWait(progressBar.Dismiss)
}
