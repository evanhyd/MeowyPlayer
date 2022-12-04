package seeker

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/resource"
)

var TheUniquePlayer *MusicPlayer

const (
	//music quality
	SAMPLING_RATE   = 44100
	NUM_OF_CHANNELS = 2
	AUDIO_BIT_DEPTH = 2

	LOOP_ORDER   = 0
	REPEAT_ORDER = 1
	RANDOM_ORDER = 2
	ORDER_LEN    = 3
)

type MusicPlayer struct {
	UpdateMusicInfo chan custom_canvas.MusicInfo
	UpdateProgress  chan float64
	UpdatePlay      chan bool
	UpdateOrder     chan int

	RequestPrev     chan struct{}
	RequestNext     chan struct{}
	RequestAdhoc    chan struct{}
	RequestPlay     chan struct{}
	RequestOrder    chan struct{}
	RequestVolume   chan float64
	RequestProgress chan float64

	otoCtx         *oto.Context
	musicInfoList  []custom_canvas.MusicInfo
	currMusicIndex int
}

func NewMusicPlayer() (*MusicPlayer, error) {
	p := &MusicPlayer{}
	p.UpdateMusicInfo = make(chan custom_canvas.MusicInfo, 32)
	p.UpdateProgress = make(chan float64, 32)
	p.UpdatePlay = make(chan bool, 32)
	p.UpdateOrder = make(chan int, 32)

	p.RequestPrev = make(chan struct{}, 32)
	p.RequestNext = make(chan struct{}, 32)
	p.RequestAdhoc = make(chan struct{}, 32)
	p.RequestPlay = make(chan struct{}, 32)
	p.RequestOrder = make(chan struct{}, 32)
	p.RequestVolume = make(chan float64, 32)
	p.RequestProgress = make(chan float64, 32)

	var ready chan struct{}
	var err error
	p.otoCtx, ready, err = oto.NewContext(SAMPLING_RATE, NUM_OF_CHANNELS, AUDIO_BIT_DEPTH)
	if err != nil {
		return nil, err
	}
	<-ready

	p.musicInfoList = make([]custom_canvas.MusicInfo, 0)
	p.currMusicIndex = 0

	return p, nil
}

/**
Set the music queue to the provided music info.
Reload and play the music given by the index.
*/
func (p *MusicPlayer) SetPlaylist(musicInfo []custom_canvas.MusicInfo, musicIndex int) {
	p.musicInfoList = musicInfo
	p.currMusicIndex = musicIndex
	p.RequestAdhoc <- struct{}{}
}

/**
Go to the previous song determined by the playing order.
*/
func (p *MusicPlayer) goPrev(order int) {

	switch order {
	case LOOP_ORDER:
		p.currMusicIndex--
		if p.currMusicIndex < 0 {
			p.currMusicIndex += len(p.musicInfoList)
		}

	case REPEAT_ORDER:

	case RANDOM_ORDER:
		for {
			newIndex := rand.Int() % len(p.musicInfoList)
			if newIndex != p.currMusicIndex {
				p.currMusicIndex = newIndex
				break
			}
		}
	}
}

/**
Go to the next song determined by the playing order.
*/
func (p *MusicPlayer) goNext(order int) {
	switch order {
	case LOOP_ORDER:
		p.currMusicIndex++
		if p.currMusicIndex >= len(p.musicInfoList) {
			p.currMusicIndex -= len(p.musicInfoList)
		}

	case REPEAT_ORDER:

	case RANDOM_ORDER:
		for {
			newIndex := rand.Int() % len(p.musicInfoList)
			if newIndex != p.currMusicIndex {
				p.currMusicIndex = newIndex
				break
			}
		}
	}
}

/**
Change the playing order to the next mode.
*/
func (p *MusicPlayer) incOrder(order int) int {
	order++
	if order >= ORDER_LEN {
		order -= ORDER_LEN
	}
	p.UpdateOrder <- order
	return order
}

/**
Launch me to start the music player!
*/
func (p *MusicPlayer) launch() {

	volume := 1.0
	order := RANDOM_ORDER
	// lock := sync.Mutex{}

	for {
		if len(p.musicInfoList) > 0 {

			//load music file
			musicInfo := p.musicInfoList[p.currMusicIndex]
			mp3File, err := os.Open(resource.GetMusicPath(musicInfo.Title))
			if err != nil {
				log.Panic(err)
			}

			//decode music file
			decodedMp3File, err := mp3.NewDecoder(mp3File)
			if err != nil {
				log.Panic(err)
			}

			//obtain music player
			player := p.otoCtx.NewPlayer(decodedMp3File)
			player.SetVolume(volume)
			p.UpdateProgress <- 0.0
			p.UpdateMusicInfo <- musicInfo
			p.UpdatePlay <- true

			player.Play()
			pausedByUser := false
			forceReload := false
			for {

				if !player.IsPlaying() && !pausedByUser || forceReload {
					break
				}

				select {

				//change playing status
				case <-p.RequestPlay:
					pausedByUser = !pausedByUser
					if pausedByUser {
						player.Pause()
					} else {
						player.Play()
					}
					p.UpdatePlay <- !pausedByUser

				//roll back song
				case <-p.RequestPrev:
					p.goPrev(order)
					forceReload = true

				//skip song
				case <-p.RequestNext:
					p.goNext(order)
					forceReload = true

				case <-p.RequestAdhoc:
					forceReload = true

				//set to user wanted playing order
				case <-p.RequestOrder:
					order = p.incOrder(order)

				//set to user wanted volume
				case newVolume := <-p.RequestVolume:
					volume = newVolume
					player.SetVolume(volume)

				//set to user wanted music progress
				case percent := <-p.RequestProgress:
					//it forces thread syncrhonization?????
					log.Println("magic begin")
					// lock.Lock()
					player.Pause()
					tick := int64(percent * float64(decodedMp3File.Length()))
					tick -= tick % 4
					decodedMp3File.Seek(tick, io.SeekStart)
					player.Play()
					p.UpdatePlay <- true
					// lock.Unlock()
					log.Println("magic end")

				//update music progress
				default:
					currTick, err := decodedMp3File.Seek(0, io.SeekCurrent)
					if err != nil {
						log.Panic(err)
					}
					percent := float64(currTick) / float64(decodedMp3File.Length())
					p.UpdateProgress <- percent
					time.Sleep(100.0 * time.Millisecond)
				}
			}

			//normal playing sequence
			if !forceReload {
				p.goNext(order)
			}

			//close the resources
			err = player.Close()
			if err != nil {
				log.Panic(err)
			}
			err = mp3File.Close()
			if err != nil {
				log.Panic(err)
			}

		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
