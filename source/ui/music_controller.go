package ui

import (
	stdcontext "context"
	"fmt"
	"log/slog"
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

type MusicController struct {
	widget.BaseWidget
	userContext *storages.UserContext
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

func newMusicController(userContext *storages.UserContext) *MusicController {
	c := &MusicController{
		userContext:    userContext,
		playlistCover:  canvas.NewImageFromResource(resourceIconPng),
		title:          widget.NewRichText(&widget.TextSegment{Style: widget.RichTextStyle{SizeName: theme.SizeNameSubHeadingText, TextStyle: fyne.TextStyle{Bold: true}}}),
		durationLabel:  widget.NewLabel("00:00"),
		modeDropDown:   mwidget.NewDropDown(),
		skipPrevButton: widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil),
		playButton:     widget.NewButtonWithIcon("", theme.MediaPauseIcon(), nil),
		skipNextButton: widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil),
		volumeSlider:   mwidget.NewVolumeSlider(),
	}

	c.musicPlayer = players.MakeBeepPlayer(userContext, c.downloadMusic, c.onPlayMusic)

	c.playlistCover.SetMinSize(mwidget.ControllerPlaylistCardSize)
	c.playlistCover.CornerRadius = 8.0
	c.title.Truncation = fyne.TextTruncateEllipsis

	c.progressSlider = mwidget.NewProgressSlider(func(percent float64) {
		c.musicPlayer.SetProgress(percent)
		c.playButton.SetIcon(theme.MediaPauseIcon())
	})

	c.modeDropDown.Add(fyne.NewMenuItemWithIcon(lang.L("Sequential"), theme.MailForwardIcon(), func() { c.musicPlayer.SetQueueMode(players.SequentialQueueMode) }))
	c.modeDropDown.Add(fyne.NewMenuItemWithIcon(lang.L("Random"), resourceRandomSvg, func() { c.musicPlayer.SetQueueMode(players.RandomQueueMode) }))
	c.modeDropDown.Select(0)

	c.skipPrevButton.Importance = widget.LowImportance
	c.skipPrevButton.OnTapped = func() {
		c.musicPlayer.Previous()
		c.playButton.SetIcon(theme.MediaPauseIcon())
	}

	c.playButton.Importance = widget.LowImportance
	c.playButton.OnTapped = func() {
		if c.musicPlayer.IsPlaying() {
			c.musicPlayer.Pause()
			c.playButton.SetIcon(theme.MediaPlayIcon())
		} else {
			c.musicPlayer.Resume()
			c.playButton.SetIcon(theme.MediaPauseIcon())
		}
	}

	c.skipNextButton.Importance = widget.LowImportance
	c.skipNextButton.OnTapped = func() {
		c.musicPlayer.Next()
		c.playButton.SetIcon(theme.MediaPauseIcon())
	}

	c.volumeSlider.OnChanged = c.musicPlayer.SetVolume
	c.volumeSlider.SetVolume(0.7)

	c.userContext.AddListener(storages.OnPlayMusicEvent, func(data any) {
		evt := data.(storages.OnPlayMusicEventData)
		c.loadPlaylist(evt.Playlist, evt.Music)
	})

	go func() {
		for range time.NewTicker(time.Second).C {
			progress := c.musicPlayer.GetProgress()
			music := c.musicPlayer.GetMusic()
			playedDuration := int64(float64(music.LengthSeconds) * progress)

			fyne.DoAndWait(func() {
				c.durationLabel.SetText(fmt.Sprintf("%s / %s", mutil.SecondsToTime(playedDuration), mutil.SecondsToTime(music.LengthSeconds)))
				c.progressSlider.SetValue(progress)
			})
		}
	}()

	c.ExtendBaseWidget(c)
	return c
}

func (c *MusicController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil, nil, c.playlistCover, nil,
		container.NewBorder(
			c.title,
			mcontainer.NewHSplit(0.8, container.NewCenter(container.NewHBox(c.modeDropDown, c.skipPrevButton, c.playButton, c.skipNextButton)), c.volumeSlider),
			nil, nil,
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
	if err != nil || m.Picture() == nil {
		return nil
	}
	return m.Picture().Data
}

func (c *MusicController) onPlayMusic(music storages.Music) {
	fyne.DoAndWait(func() {
		c.title.Segments[0].(*widget.TextSegment).Text = music.Title
		c.title.Refresh()

		if cover := c.extractCover(music); cover != nil {
			c.playlistCover.Resource = fyne.NewStaticResource(music.MusicId, cover)
		} else {
			c.playlistCover.Resource = fyne.NewStaticResource(music.MusicId, c.playlist.CoverBlob)
		}
		c.playlistCover.Refresh()
	})
}

func (c *MusicController) loadPlaylist(playlist storages.Playlist, music storages.Music) {
	c.playlist = playlist
	c.musicPlayer.SetPlaylist(playlist, music)
	c.playButton.SetIcon(theme.MediaPauseIcon())
}

func (c *MusicController) downloadMusic(music storages.Music) {
	var progressBar dialog.Dialog

	fyne.DoAndWait(func() {
		progressBar = dialog.NewCustomWithoutButtons(lang.L("Downloading"), widget.NewProgressBarInfinite(), fyne.CurrentApp().Driver().AllWindows()[0])
		progressBar.Show()
	})
	defer fyne.DoAndWait(progressBar.Dismiss)

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

	if err = c.userContext.PutMusicFile(music, content); err != nil {
		slog.Error("failed to put music file", "error", err)
	}
}
