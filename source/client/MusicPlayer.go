package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/container"
)

const (
	RANDOM = iota
	ORDERED
	REPEAT
	SIZE
)

type MusicPlayer struct {
	resource.PlayList
	playMode    int
	history     container.Vector[int]
	randomQueue container.Vector[int]

	//channel to syncrhonize the commands
	playListChan chan resource.PlayList
	progressCMD  chan float64
	volumeCMD    chan float64
	playCMD      chan struct{}
	rollbackCMD  chan struct{}
	skipCMD      chan struct{}
	modeCMD      chan int
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		playMode:     RANDOM,
		playListChan: make(chan resource.PlayList),
		progressCMD:  make(chan float64, 16),
		volumeCMD:    make(chan float64, 16),
		playCMD:      make(chan struct{}, 16),
		rollbackCMD:  make(chan struct{}, 16),
		skipCMD:      make(chan struct{}, 16),
		modeCMD:      make(chan int, 16),
	}
}

func (m *MusicPlayer) CommandProgress(percent float64) {
	m.progressCMD <- percent
}

func (m *MusicPlayer) CommandVolume(volume float64) {
	m.volumeCMD <- volume
}

func (m *MusicPlayer) CommandPlay() {
	m.playCMD <- struct{}{}
}

func (m *MusicPlayer) CommandRollback() {
	m.rollbackCMD <- struct{}{}
}

func (m *MusicPlayer) CommandSkip() {
	m.skipCMD <- struct{}{}
}

func (m *MusicPlayer) CommandMode(mode int) {
	m.modeCMD <- mode
}

func (m *MusicPlayer) setPlayMode(playMode int) {
	if playMode == RANDOM {
		m.history.Clear()
		m.randomQueue = rand.Perm(len(m.Album().MusicList))
	}
	m.playMode = playMode
}

func (m *MusicPlayer) setPlayList(playList resource.PlayList) {
	m.PlayList = playList
	m.setPlayMode(m.playMode)
}

func (m *MusicPlayer) rollback() {
	switch m.playMode {
	case RANDOM:
		if !m.history.Empty() {
			m.randomQueue.PushBack(m.Index())
			m.SetIndex(*m.history.Back())
			m.history.PopBack()
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
		if m.randomQueue.Empty() {
			m.randomQueue = rand.Perm(len(m.Album().MusicList))
		}
		m.history.PushBack(m.Index())
		m.SetIndex(*m.randomQueue.Back())
		m.randomQueue.PopBack()

	case ORDERED:
		m.SetIndex((m.Index() + 1) % len(m.Album().MusicList))

	case REPEAT:
		//nothing
	}
}

func (m *MusicPlayer) Notify(play *resource.PlayList) {
	m.playListChan <- *play
}

func (m *MusicPlayer) Start(menu *cwidget.CommandMenu) {
	//initialize oto mp3 context
	context, ready, err := oto.NewContext(resource.SAMPLING_RATE, resource.NUM_OF_CHANNELS, resource.AUDIO_BIT_DEPTH)
	assert.NoErr(err, "failed to start the music player")
	<-ready

	//wait for the user to click the music
WaitLoop:
	for {
		select {
		case playList := <-m.playListChan:
			m.setPlayList(playList)
			break WaitLoop
		case <-m.skipCMD:
		case <-m.rollbackCMD:
		case <-m.playCMD:
		case <-m.modeCMD:
		case <-m.progressCMD:
		case <-m.volumeCMD:
			//drain out meaningless commands
		}
	}

	for {
		menu.SetMusic(m.Music())
		mp3Controller := resource.NewMP3Controller(context, m.PlayList.Music())
		interrupted := false
		fmt.Println("playing", m.Music().SimpleTitle())

	CONTROL_LOOP:
		for mp3Controller.PlayOrPause(); !mp3Controller.IsOver(); {
			select {
			case playList := <-m.playListChan:
				fmt.Println("new playlist")
				m.setPlayList(playList)
				interrupted = true
				break CONTROL_LOOP

			case <-m.skipCMD:
				fmt.Println("skip")
				m.skip()
				interrupted = true
				break CONTROL_LOOP

			case <-m.rollbackCMD:
				fmt.Println("rollback")
				m.rollback()
				interrupted = true
				break CONTROL_LOOP

			case <-m.playCMD:
				fmt.Println("play/pause")
				mp3Controller.PlayOrPause()

			case playMode := <-m.modeCMD:
				fmt.Println("mode", playMode)
				m.setPlayMode(playMode)

			case percent := <-m.progressCMD:
				fmt.Println("progress", percent)
				mp3Controller.Pause()
				time.Sleep(10 * time.Millisecond) //mp3 library has race condition
				mp3Controller.SetProgress(percent)
				mp3Controller.Play()

			case volume := <-m.volumeCMD:
				fmt.Println("set volume")
				mp3Controller.SetVolume(volume)

			default:
				time.Sleep(100 * time.Millisecond)
			}
			menu.UpdateProgress(m.PlayList.Music().Length, mp3Controller.CurrentProgressPercent())
		}

		if !interrupted {
			m.skip()
		}

		mp3Controller.Close()
	}
}
