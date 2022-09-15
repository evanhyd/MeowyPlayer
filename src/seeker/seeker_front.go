package seeker

import (
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/src/resource"
)

//music controller
func NewSeekerUI() *fyne.Container {

	var err error
	MusicQueue, err = NewMusicSeeker()
	if err != nil {
		log.Fatal(err)
	}

	//launch seeker
	go LaunchSeeker()

	bot := container.NewVBox(
		container.NewMax(layout.NewSpacer(), MusicQueue.Progress, layout.NewSpacer()),
		container.NewBorder(nil, nil, MusicQueue.Title, container.NewHBox(MusicQueue.Prev, MusicQueue.Play, MusicQueue.Next, MusicQueue.Order, MusicQueue.Volume)))
	return container.NewBorder(nil, bot, nil, nil)
}

func LaunchSeeker() {

	//create oto context
	log.Println("loading driver...")
	otoCtx, ready, err := oto.NewContext(SAMPLING_RATE, NUM_OF_CHANNELS, AUDIO_BIT_DEPTH)
	if err != nil {
		log.Panic(err)
	}

	//wait for the hardware to get ready
	<-ready
	log.Println("driver is ready")

	for {
		if !MusicQueue.IsEmpty() {

			//load music file
			info := MusicQueue.GetCurrMusic()
			musicFile, err := os.Open(resource.GetMusicPath(info.Title))
			if err != nil {
				log.Panic(err)
			}

			//decode music file
			decodedMP3, err := mp3.NewDecoder(musicFile)
			if err != nil {
				log.Panic(err)
			}

			//obtain music player
			player := otoCtx.NewPlayer(decodedMP3)

			//connect music progress bar
			MusicQueue.Progress.BindMP3(decodedMP3, &player)

			//connect volume bar
			MusicQueue.Volume.BindMP3(decodedMP3, &player)

			//set up title
			MusicQueue.Title.SetText(info.Title)
			log.Println("playing: " + info.Title)

			shouldPlay := true
			shouldReload := false
			player.Play()
			MusicQueue.Play.UpdateIcon(true)
			for (player.IsPlaying() || !shouldPlay) && !shouldReload {

				select {
				//play or pause
				case <-MusicQueue.Play.Signal:
					shouldPlay = !shouldPlay
					if shouldPlay {
						player.Play()
						log.Println("resumed")
					} else {
						player.Pause()
						log.Println("paused")
					}
					MusicQueue.Play.UpdateIcon(shouldPlay)

					//rewind to previous song
				case <-MusicQueue.Prev.Signal:
					shouldReload = true
					log.Println("prev")

					//skip to next song
				case <-MusicQueue.Next.Signal:
					shouldReload = true
					log.Println("next")

				default:

					//update the music progress bar
					if shouldPlay {
						MusicQueue.Progress.UpdateBar()
					}
					time.Sleep(time.Second)
				}
			}

			//check if should proceed normally
			if !shouldReload {
				MusicQueue.NextMusic()
			}

			err = player.Close()
			if err != nil {
				log.Panic(err)
			}
			err = musicFile.Close()
			if err != nil {
				log.Panic(err)
			}

		} else {
			log.Println("empty queue...")
			time.Sleep(time.Second)
		}
	}
}
