package seeker

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/src/resource"
)

//music controller
func NewSeekerUI() *fyne.Container {

	log.Println("loading music player...")
	var err error
	TheUniquePlayer, err = NewMusicPlayer()
	if err != nil {
		log.Println(err)
	}

	log.Println("loading music player textures")
	playingIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("seeker_play.png"))
	if err != nil {
		log.Println(err)
	}
	pausingIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("seeker_pause.png"))
	if err != nil {
		log.Println(err)
	}
	modeIcons := make([]fyne.Resource, 0, MODE_LEN)
	for i := 0; i < MODE_LEN; i++ {
		icon, err := fyne.LoadResourceFromPath(resource.GetImagePath(fmt.Sprintf("seeker_mode_icon_%v.png", i)))
		if err != nil {
			log.Panic(err)
		}
		modeIcons = append(modeIcons, icon)
	}

	title := widget.NewLabel("")
	progressTitle := widget.NewLabel("00.00%")
	progressBar := widget.NewSlider(0.0, 1.0)
	progressBar.Step = 0.00000001
	progressBar.OnChanged = func(val float64) { TheUniquePlayer.RequestProgress <- val }
	prevBtn := widget.NewButton("<<", func() { TheUniquePlayer.RequestPrev <- struct{}{} })
	nextBtn := widget.NewButton(">>", func() { TheUniquePlayer.RequestNext <- struct{}{} })
	playBtn := widget.NewButtonWithIcon("", pausingIcon, func() { TheUniquePlayer.RequestPlay <- struct{}{} })
	modeBtn := widget.NewButtonWithIcon("", modeIcons[RANDOM_MODE], func() { TheUniquePlayer.RequestMode <- struct{}{} })
	volume := widget.NewSlider(0.0, 1.0)
	volume.SetValue(1.0)
	volume.Step = 0.01
	volume.OnChanged = func(val float64) { TheUniquePlayer.RequestVolume <- val }

	go func() {
		for {
			select {
			//update music title
			case musicInfo := <-TheUniquePlayer.UpdateMusicTitle:
				title.SetText(musicInfo.Title)
				log.Println("playing: " + musicInfo.Title)

			//update music progress
			case percent := <-TheUniquePlayer.UpdateProgressBar:
				progressTitle.SetText(fmt.Sprintf("%05.2f%%", percent*100))

				//avoid SetValue triggering OnChanged()
				progressBar.Value = percent
				progressBar.Refresh()

			//update play button
			case isPlaying := <-TheUniquePlayer.UpdatePlayIcon:
				if isPlaying {
					playBtn.SetIcon(pausingIcon)
				} else {
					playBtn.SetIcon(playingIcon)
				}

			//update play mode
			case mode := <-TheUniquePlayer.UpdateModeIcon:
				modeBtn.SetIcon(modeIcons[mode])
			}
		}
	}()

	go TheUniquePlayer.Launch()

	return container.NewBorder(
		title, nil, nil, nil,
		container.NewVBox(
			container.NewBorder(nil, nil, progressTitle, nil, progressBar),
			container.NewHBox(layout.NewSpacer(), prevBtn, playBtn, nextBtn, modeBtn, volume, layout.NewSpacer()),
		),
	)
}
