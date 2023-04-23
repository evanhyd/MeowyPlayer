package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/cwidget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

var seekerPlayModeIcons [player.PLAYMODE_LEN]fyne.Resource

func init() {
	const (
		seekerPlayModeRandomIconName  = "seeker_playmode_random.png"
		seekerPlayModeOrderedIconName = "seeker_playmode_ordered.png"
		seekerPlayModeRepeatIconName  = "seeker_playmode_repeat.png"
	)

	var err error
	if seekerPlayModeIcons[player.RANDOM], err = fyne.LoadResourceFromPath(resource.GetResourcePath(seekerPlayModeRandomIconName)); err != nil {
		log.Fatal(err)
	}
	if seekerPlayModeIcons[player.ORDERED], err = fyne.LoadResourceFromPath(resource.GetResourcePath(seekerPlayModeOrderedIconName)); err != nil {
		log.Fatal(err)
	}
	if seekerPlayModeIcons[player.REPEAT], err = fyne.LoadResourceFromPath(resource.GetResourcePath(seekerPlayModeRepeatIconName)); err != nil {
		log.Fatal(err)
	}
}

func createSeeker() *fyne.Container {
	albumView := cwidget.NewCardWithImage("", "", nil, nil)
	albumView.Image = canvas.NewImageFromResource(defaultIcon)
	albumView.Image.SetMinSize(resource.GetAlbumCoverSize())
	player.GetState().OnUpdateSeeker().AddCallback(func(album player.Album, _ []player.Music, _ player.Music) {
		albumView.SetImage(album.CoverIcon())
		albumView.OnTapped = func() { DisplayErrorIfAny(player.UserSelectAlbum(album)) }
	})

	title := widget.NewLabel("")
	player.GetPlayer().OnMusicBeginSubject().AddCallback(func(music player.Music) { title.SetText(music.Title()[:len(music.Title())-4]) })

	progressLabel := widget.NewLabel("00:00")
	player.GetPlayer().OnMusicPlayingSubject().AddCallback(func(music player.Music, percent float64) {
		secPassed := int(music.Duration().Seconds() * percent)
		progressLabel.SetText(fmt.Sprintf("%02d:%02d", secPassed/60, secPassed%60))
	})

	progress := widget.NewSlider(0.0, 1.0)
	progress.Step = 1.0 / 10000.0
	progress.OnChanged = func(percent float64) { player.GetPlayer().SetProgress(percent) }
	player.GetPlayer().OnMusicPlayingSubject().AddCallback(func(music player.Music, percent float64) {
		progress.Value = percent
		progress.Refresh()
	})

	prevButton := cwidget.NewButton(" << ")
	playButton := cwidget.NewButton(" O ")
	nextButton := cwidget.NewButton(" >> ")
	prevButton.OnTapped = player.GetPlayer().PreviousMusic
	playButton.OnTapped = player.GetPlayer().PlayPauseMusic
	nextButton.OnTapped = player.GetPlayer().NextMusic

	playMode := player.RANDOM
	playModeButton := cwidget.NewButtonWithIcon("", seekerPlayModeIcons[playMode])
	playModeButton.OnTapped = func() {
		playMode = (playMode + 1) % player.PLAYMODE_LEN
		playModeButton.SetIcon(seekerPlayModeIcons[playMode])
		player.GetPlayer().SetPlayMode(playMode)
	}

	volume := widget.NewSlider(0.0, 1.0)
	volume.Step = 1.0 / 100.0
	volume.Value = 1.0
	volume.OnChanged = func(volume float64) { player.GetPlayer().SetMusicVolume(volume) }

	return container.NewBorder(
		nil,
		nil,
		albumView,
		nil,
		container.NewBorder(
			title,
			container.NewHBox(layout.NewSpacer(), prevButton, playButton, nextButton, playModeButton, volume, layout.NewSpacer()),
			nil,
			nil,
			container.NewBorder(nil, nil, progressLabel, nil, progress),
		),
	)
}
