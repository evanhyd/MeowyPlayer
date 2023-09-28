package client

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/source/utility"
)

type MusicPlayer struct {
	*oto.Context
	playListChan chan player.PlayList
}

func NewMusicPlayer() *MusicPlayer {
	//initialize oto mp3 player
	context, ready, err := oto.NewContext(player.SAMPLING_RATE, player.NUM_OF_CHANNELS, player.AUDIO_BIT_DEPTH)
	utility.MustNil(err)
	<-ready

	return &MusicPlayer{
		Context:      context,
		playListChan: make(chan player.PlayList),
	}
}

func (m *MusicPlayer) Notify(play *player.PlayList) {
	m.playListChan <- *play
}

func (m *MusicPlayer) decode(music *player.Music) MP3Controller {
	mp3Data, err := os.ReadFile(resource.MusicPath(music))
	utility.MustNil(err)
	mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	utility.MustNil(err)
	return makeMP3Player(mp3Decoder, m.NewPlayer(mp3Decoder))
}

func (m *MusicPlayer) Start(menu *cwidget.PlayerMenu) {
	const (
		NORMAL = iota
		SKIP
		ROLLBACK
	)
	menuChannel := menu.GetMenuChannel()

	//wait for the user to click the music
	playList := <-m.playListChan
	for {
		log.Printf("playing %v\n", playList.Music().SimpleTitle())
		menu.SetMusic(playList.Music())
		mp3Controller := m.decode(playList.Music())
		completeStatus := NORMAL

	CONTROL_LOOP:
		for mp3Controller.Play(); !mp3Controller.IsOver(); {
			select {
			case playList = <-m.playListChan:
				fmt.Println("new playlist")
				break CONTROL_LOOP

			case <-menuChannel.Skip:
				fmt.Println("skip")
				completeStatus = SKIP
				break CONTROL_LOOP

			case <-menuChannel.Rollback:
				fmt.Println("rollback")
				completeStatus = ROLLBACK
				break CONTROL_LOOP

			case <-menuChannel.Play:
				fmt.Println("play/pause")
				if mp3Controller.IsPlaying() {
					mp3Controller.Pause()
				} else {
					mp3Controller.Play()
				}

			case percent := <-menuChannel.Progress:
				fmt.Println("progress", percent)
				mp3Controller.Pause()
				time.Sleep(10 * time.Millisecond) //mp3 library has race condition
				mp3Controller.SetProgress(percent)
				mp3Controller.Play()

			case volume := <-menuChannel.Volume:
				fmt.Println("set volume")
				mp3Controller.SetVolume(volume)

			default:
				menu.UpdateProgressBar(mp3Controller.CurrentProgressPercent())
				time.Sleep(100 * time.Millisecond)
			}
		}

		mp3Controller.Close()

		switch completeStatus {
		case NORMAL, SKIP:
			playList.NextMusic()
		case ROLLBACK:
			playList.PrevMusic()
		}
	}
}
