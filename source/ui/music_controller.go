package ui

import (
	"meowyplayer/context"
	"meowyplayer/players"
	"meowyplayer/storages"
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mutil"
	"meowyplayer/ui/internal/mwidget"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	Pause  = 0
	Resume = 1
)

type MusicController struct {
	widget.BaseWidget
	userContext *context.UserContext
	musicPlayer players.MusicPlayer

	playlistCover  *canvas.Image
	title          *widget.RichText
	progressSlider *mwidget.ProgressSlider
	durationLabel  *widget.Label
	modeDropDown   *mwidget.DropDown
	skipPrevButton *widget.Button
	playButton     *mwidget.MultiButton
	skipNextButton *widget.Button
	volumeSlider   *widget.Slider
}

func newMusicController(userContext *context.UserContext) *MusicController {
	var c MusicController
	c = MusicController{
		userContext:   userContext,
		musicPlayer:   players.MakeBeepPlayer(userContext),
		playlistCover: canvas.NewImageFromResource(resourceIconPng),
		title: widget.NewRichText(&widget.TextSegment{
			Style: widget.RichTextStyle{SizeName: theme.SizeNameHeadingText, TextStyle: fyne.TextStyle{Bold: true}},
		}),
		durationLabel:  widget.NewLabel("00:00"),
		modeDropDown:   mwidget.NewDropDown(),
		skipPrevButton: widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil),
		playButton:     mwidget.NewMultiButton(),
		skipNextButton: widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil),
		volumeSlider:   widget.NewSlider(0.0, 1.0),
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
	c.playButton.Add(theme.MediaPlayIcon(), c.musicPlayer.Pause)
	c.playButton.Add(theme.MediaPauseIcon(), c.musicPlayer.Resume)
	c.playButton.Select(0)

	c.skipNextButton.Importance = widget.LowImportance
	c.skipNextButton.OnTapped = c.musicPlayer.Next

	c.volumeSlider.Step = 0.01
	c.volumeSlider.OnChanged = c.musicPlayer.SetVolume

	// Some music player initialization.
	c.volumeSlider.SetValue(0.7)

	c.userContext.AddListener(context.OnPlayPlaylistEvent, func(_ context.EventType, data any) {
		c.fetchMusic(data.(context.OnPlayPlaylistEventData).Playlist, data.(context.OnPlayPlaylistEventData).SelectedMusic)
	})

	// UI update thread
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			progress := c.musicPlayer.GetProgress()
			music := c.musicPlayer.GetMusic()
			remainDuration := music.LengthSeconds - int64(float64(music.LengthSeconds)*progress)
			fyne.DoAndWait(func() {
				c.durationLabel.SetText(mutil.SecondsToTime(remainDuration))
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

func (c *MusicController) fetchMusic(playlist storages.Playlist, selectedMusic storages.Music) {
	go c.musicPlayer.SetPlaylist(playlist, selectedMusic)
	c.playlistCover.Resource = fyne.NewStaticResource(mutil.PlaylistIdToString(playlist.PlaylistId), playlist.CoverBlob)
	c.playlistCover.Refresh()
}
