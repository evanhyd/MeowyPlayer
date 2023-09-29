package client

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/source/utility"
)

type MusicOrder int

const (
	RANDOM MusicOrder = iota
	ORDERED
	REPEAT
	SIZE
)

type MusicPlayer struct {
	*oto.Context
	player.PlayList
	playListChan chan player.PlayList

	playMode    MusicOrder
	history     []int
	randomQueue []int
}

func NewMusicPlayer() *MusicPlayer {
	//initialize oto mp3 player
	context, ready, err := oto.NewContext(player.SAMPLING_RATE, player.NUM_OF_CHANNELS, player.AUDIO_BIT_DEPTH)
	utility.MustNil(err)
	<-ready

	return &MusicPlayer{Context: context, playListChan: make(chan player.PlayList), playMode: RANDOM}
}

func (m *MusicPlayer) setPlayMode(playMode MusicOrder) {
	if playMode == RANDOM {
		m.history = []int{}
		m.randomQueue = rand.Perm(len(m.Album().MusicList))
	}
	m.playMode = playMode
}

func (m *MusicPlayer) setPlayList(playList player.PlayList) {
	m.PlayList = playList
	m.setPlayMode(m.playMode)
}

func (m *MusicPlayer) decode(music *player.Music) MP3Controller {
	mp3Data, err := os.ReadFile(resource.MusicPath(music))
	utility.MustNil(err)
	mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	utility.MustNil(err)
	return makeMP3Player(mp3Decoder, m.NewPlayer(mp3Decoder))
}

func (m *MusicPlayer) rollback() {
	switch m.playMode {
	case RANDOM:
		if len(m.history) > 0 {
			m.randomQueue = append(m.randomQueue, m.Index())
			last := len(m.history) - 1
			m.SetIndex(m.history[last])
			m.history = m.history[:last]
		}

	case ORDERED:
		m.SetIndex((m.Index() - 1 + len(m.Album().MusicList)) % len(m.Album().MusicList))

	case REPEAT:
		//nothing
	}
}

func (m *MusicPlayer) skip() {
	switch m.playMode {
	case RANDOM:
		//generate new queue if run out of music
		if len(m.randomQueue) == 0 {
			m.randomQueue = rand.Perm(len(m.Album().MusicList))
		}
		m.history = append(m.history, m.Index())
		last := len(m.randomQueue) - 1
		m.SetIndex(m.randomQueue[last])
		m.randomQueue = m.randomQueue[:last]

	case ORDERED:
		m.SetIndex((m.Index() + 1) % len(m.Album().MusicList))

	case REPEAT:
		//nothing
	}
}

func (m *MusicPlayer) Notify(play *player.PlayList) {
	m.playListChan <- *play
}

func (m *MusicPlayer) Start(menu *cwidget.PlayerMenu) {
	menuChannel := menu.GetMenuChannel()

	//wait for the user to click the music
	m.setPlayList(<-m.playListChan)
	for {
		log.Printf("playing %v\n", m.Music().SimpleTitle())
		menu.SetMusic(m.Music())
		mp3Controller := m.decode(m.Music())
		interrupted := false

	CONTROL_LOOP:
		for mp3Controller.Play(); !mp3Controller.IsOver(); {
			select {
			case playList := <-m.playListChan:
				fmt.Println("new playlist")
				m.setPlayList(playList)
				interrupted = true
				break CONTROL_LOOP

			case <-menuChannel.Skip:
				fmt.Println("skip")
				m.skip()
				interrupted = true
				break CONTROL_LOOP

			case <-menuChannel.Rollback:
				fmt.Println("rollback")
				m.rollback()
				interrupted = true
				break CONTROL_LOOP

			case <-menuChannel.Play:
				fmt.Println("play/pause")
				mp3Controller.PlayOrPause()

			case <-menuChannel.Mode:
				m.setPlayMode((m.playMode + 1) % SIZE)

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
			}
			menu.UpdateProgressBar(mp3Controller.CurrentProgressPercent())
			time.Sleep(100 * time.Millisecond)
		}

		if !interrupted {
			m.skip()
		}

		mp3Controller.Close()
	}
}
