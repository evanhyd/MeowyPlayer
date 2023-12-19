package player

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/container"
	"meowyplayer.com/utility/logger"
)

const (
	RANDOM = iota
	ORDERED
	REPLAY
	SIZE
)

type MP3Player struct {
	PlayList
	playMode    int
	history     container.Slice[int]
	randomQueue container.Slice[int]

	//channel to syncrhonize the commands
	playListChan chan PlayList
	progressCMD  chan float64
	volumeCMD    chan float64
	playCMD      chan struct{}
	rollbackCMD  chan struct{}
	skipCMD      chan struct{}
	modeCMD      chan int
}

func NewMP3Player() *MP3Player {
	return &MP3Player{
		playMode:     RANDOM,
		playListChan: make(chan PlayList),
		progressCMD:  make(chan float64, 16),
		volumeCMD:    make(chan float64, 16),
		playCMD:      make(chan struct{}, 16),
		rollbackCMD:  make(chan struct{}, 16),
		skipCMD:      make(chan struct{}, 16),
		modeCMD:      make(chan int, 16),
	}
}

func (m *MP3Player) OnProgress(percent float64) {
	m.progressCMD <- percent
}

func (m *MP3Player) OnVolume(volume float64) {
	m.volumeCMD <- volume
}

func (m *MP3Player) OnPlay() {
	m.playCMD <- struct{}{}
}

func (m *MP3Player) OnRollback() {
	m.rollbackCMD <- struct{}{}
}

func (m *MP3Player) OnSkip() {
	m.skipCMD <- struct{}{}
}

func (m *MP3Player) OnMode(mode int) {
	m.modeCMD <- mode
}

func (m *MP3Player) setPlayMode(playMode int) {
	if playMode == RANDOM {
		m.history.Clear()
		m.randomQueue = rand.Perm(m.MusicCount())
	}
	m.playMode = playMode
}

func (m *MP3Player) setPlayList(playList PlayList) {
	m.PlayList = playList
	m.setPlayMode(m.playMode)
}

func (m *MP3Player) rollback() {
	switch m.playMode {
	case RANDOM:
		if !m.history.Empty() {
			m.randomQueue.PushBack(m.Index())
			m.SetIndex(*m.history.Back())
			m.history.PopBack()
		}

	case ORDERED:
		m.SetIndex((m.Index() - 1 + m.MusicCount()) % m.MusicCount())

	case REPLAY:
		//nothing
	}
}

func (m *MP3Player) skip() {
	switch m.playMode {
	case RANDOM:
		//generate new queue if run out of music
		if m.randomQueue.Empty() {
			m.randomQueue = rand.Perm(m.MusicCount())
		}
		m.history.PushBack(m.Index())
		m.SetIndex(*m.randomQueue.Back())
		m.randomQueue.PopBack()

	case ORDERED:
		m.SetIndex((m.Index() + 1) % m.MusicCount())

	case REPLAY:
		//nothing
	}
}

func (m *MP3Player) Notify(play PlayList) {
	m.playListChan <- play
}

func (m *MP3Player) Start(menu *cwidget.MusicController) {
	//initialize oto mp3 context
	context, ready, err := oto.NewContext(resource.SAMPLING_RATE, resource.NUM_OF_CHANNELS, resource.AUDIO_BIT_DEPTH)
	if err != nil {
		logger.Error(err, 0)
	}
	<-ready

	//initialize loop timer
	idleTimer := time.Tick(1 * time.Second)

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
		mp3Controller := NewSeeker(context, m.PlayList.Music())
		mp3Controller.SetVolume(menu.Volume())

		interrupted := false

	CONTROL_LOOP:
		for mp3Controller.PlayOrPause(); !mp3Controller.IsOver(); {
			select {
			case playList := <-m.playListChan:
				m.setPlayList(playList)
				interrupted = true
				break CONTROL_LOOP

			case <-m.skipCMD:
				m.skip()
				interrupted = true
				break CONTROL_LOOP

			case <-m.rollbackCMD:
				m.rollback()
				interrupted = true
				break CONTROL_LOOP

			case <-m.playCMD:
				mp3Controller.PlayOrPause()

			case playMode := <-m.modeCMD:
				m.setPlayMode(playMode)

			case percent := <-m.progressCMD:
				mp3Controller.Pause()
				time.Sleep(10 * time.Millisecond) //mp3 library has race condition
				mp3Controller.SetProgress(percent)
				mp3Controller.Play()

			case volume := <-m.volumeCMD:
				mp3Controller.SetVolume(volume * volume) //x^2 feels more natural

			case <-idleTimer:
			}
			menu.UpdateProgress(m.PlayList.Music().Length, mp3Controller.CurrentProgressPercent())
		}

		if !interrupted {
			m.skip()
		}

		mp3Controller.Close()
	}
}
