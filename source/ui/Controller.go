package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newController() fyne.CanvasObject {

	return container.NewAdaptiveGrid(3,
		widget.NewButtonWithIcon("", theme.MenuExpandIcon(), nil),
		widget.NewButtonWithIcon("", theme.AccountIcon(), nil),
		widget.NewButtonWithIcon("", theme.ComputerIcon(), nil),
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil),
		widget.NewButtonWithIcon("", theme.HelpIcon(), nil),
		widget.NewButtonWithIcon("", theme.LoginIcon(), nil),
		widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		widget.NewButtonWithIcon("", theme.RadioButtonIcon(), nil),
		widget.NewButtonWithIcon("", theme.MediaMusicIcon(), nil),
		widget.NewButtonWithIcon("", theme.SettingsIcon(), nil),
		widget.NewButtonWithIcon("", theme.SearchIcon(), nil),
		widget.NewButtonWithIcon("", theme.UploadIcon(), nil))

	const defaultCoverName = "default.png"
	// defaultCoverSize := fyne.NewSize(128.0, 128.0)

	// coverView := cwidget.NewCardWithImage("", "", nil, nil)
	// coverView.Image = canvas.NewImageFromResource(resource.GetAsset(defaultCoverName))
	// coverView.Image.SetMinSize(defaultCoverSize)

	// manager.GetCurrentAlbum().Attach(utility.MakeCallback(func(album *player.Album) {
	// 	coverView.SetImage(canvas.NewImageFromResource(album.Cover))
	// 	coverView.OnTapped = func() { manager.GetCurrentAlbum().Set(album) }
	// }))

	// musicTitle := widget.NewLabel("title")
	// player.GetPlayer().OnMusicBeginSubject().AddCallback(func(music player.Music) { title.SetText(music.Title()[:len(music.Title())-4]) })

	// progressLabel := widget.NewLabel("00:00")
	// player.GetPlayer().OnMusicPlayingSubject().AddCallback(func(music player.Music, percent float64) {
	// 	secPassed := int(music.Duration().Seconds() * percent)
	// 	progressLabel.SetText(fmt.Sprintf("%02d:%02d", secPassed/60, secPassed%60))
	// })

	progress := widget.NewSlider(0.0, 1.0)
	progress.Step = 1.0 / 10000.0
	// progress.OnChanged = func(percent float64) { player.GetPlayer().SetProgress(percent) }
	// player.GetPlayer().OnMusicPlayingSubject().AddCallback(func(music player.Music, percent float64) {
	// 	progress.Value = percent
	// 	progress.Refresh()
	// })

	prevButton := widget.NewButton(" << ", func() {})
	playButton := widget.NewButton(" O ", func() {})
	nextButton := widget.NewButton(" >> ", func() {})
	prevButton.Importance = widget.LowImportance
	playButton.Importance = widget.LowImportance
	nextButton.Importance = widget.LowImportance
	// prevButton.OnTapped = player.GetPlayer().PreviousMusic
	// playButton.OnTapped = player.GetPlayer().PlayPauseMusic
	// nextButton.OnTapped = player.GetPlayer().NextMusic

	playModeButton := widget.NewButton("play mode", func() {})
	playModeButton.Importance = widget.LowImportance
	// playMode := player.RANDOM
	// playModeButton := cwidget.NewButtonWithIcon("", seekerPlayModeIcons[playMode])
	// playModeButton.OnTapped = func() {
	// 	playMode = (playMode + 1) % player.PLAYMODE_LEN
	// 	playModeButton.SetIcon(seekerPlayModeIcons[playMode])
	// 	player.GetPlayer().SetPlayMode(playMode)
	// }

	volume := widget.NewSlider(0.0, 1.0)
	volume.Step = 1.0 / 100.0
	volume.Value = 1.0
	// volume.OnChanged = func(volume float64) { player.GetPlayer().SetMusicVolume(volume) }
	return nil

	// return container.NewBorder(
	// 	nil,
	// 	nil,
	// 	coverView,
	// 	nil,
	// 	container.NewBorder(
	// 		musicTitle,
	// 		container.NewGridWithRows(1, layout.NewSpacer(), playModeButton, container.NewHBox(prevButton, playButton, nextButton), volume, layout.NewSpacer()),
	// 		nil,
	// 		nil,
	// 		container.NewBorder(nil, nil, progressLabel, nil, progress),
	// 	),
	// )
}
