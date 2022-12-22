package seeker

import (
	"fmt"
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

	LOOP_MODE   = 0
	REPEAT_MODE = 1
	RANDOM_MODE = 2
	MODE_LEN    = 3
)

type MusicPlayer struct {
	UpdateMusicTitle  chan custom_canvas.MusicInfo
	UpdateProgressBar chan float64
	UpdatePlayIcon    chan bool
	UpdateModeIcon    chan int

	RequestPrev     chan struct{}
	RequestNext     chan struct{}
	RequestAdhoc    chan struct{}
	RequestPlay     chan struct{}
	RequestMode     chan struct{}
	RequestVolume   chan float64
	RequestProgress chan float64

	otoCtx        *oto.Context
	musicInfoList []custom_canvas.MusicInfo
	musicIndex    int
	playMode      int
	playedMusic   map[int]struct{}
}

func NewMusicPlayer() (*MusicPlayer, error) {
	p := &MusicPlayer{}
	p.UpdateMusicTitle = make(chan custom_canvas.MusicInfo, 32)
	p.UpdateProgressBar = make(chan float64, 32)
	p.UpdatePlayIcon = make(chan bool, 32)
	p.UpdateModeIcon = make(chan int, 32)

	p.RequestPrev = make(chan struct{}, 32)
	p.RequestNext = make(chan struct{}, 32)
	p.RequestAdhoc = make(chan struct{}, 32)
	p.RequestPlay = make(chan struct{}, 32)
	p.RequestMode = make(chan struct{}, 32)
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
	p.musicIndex = 0
	p.playMode = RANDOM_MODE
	p.playedMusic = make(map[int]struct{})
	return p, nil
}

/**
Set the music queue to the provided music info.
Reload and play the music at the given index.
*/
func (p *MusicPlayer) SetPlaylist(musicInfo []custom_canvas.MusicInfo, musicIndex int) {
	p.musicInfoList = musicInfo
	p.musicIndex = musicIndex
	p.playedMusic = make(map[int]struct{})
	p.RequestAdhoc <- struct{}{}
}

/**
Go to the previous song determined by the playing mode.
*/
func (p *MusicPlayer) goPrev() {
	switch p.playMode {
	case LOOP_MODE:
		p.musicIndex = (p.musicIndex - 1 + len(p.musicInfoList)) % len(p.musicInfoList)

	case REPEAT_MODE:
		//do nothing

	case RANDOM_MODE:
	}
}

/**
Go to the next song determined by the playing mode.
*/
func (p *MusicPlayer) goNext() {
	switch p.playMode {
	case LOOP_MODE:
		p.musicIndex = (p.musicIndex + 1) % len(p.musicInfoList)

	case REPEAT_MODE:
		//do nothing

	case RANDOM_MODE:
		//all music in the list have been played
		//reset the record
		if len(p.playedMusic) == len(p.musicInfoList) {
			p.playedMusic = make(map[int]struct{})
		}

		//get a list of song indices that haven't been played
		candidate := make([]int, 0)
		for i := range p.musicInfoList {
			if _, ok := p.playedMusic[i]; !ok {
				candidate = append(candidate, i)
			}
		}
		p.musicIndex = candidate[rand.Int()%len(candidate)]
		p.playedMusic[p.musicIndex] = struct{}{}
	}
}

/**
Change the playing mode to the next mode.
*/
func (p *MusicPlayer) nextMode() {
	p.playMode = (p.playMode + 1) % MODE_LEN
	if p.playMode == RANDOM_MODE {
		p.playedMusic = make(map[int]struct{})
	}
	p.UpdateModeIcon <- p.playMode
}

/**
Get the current music info
*/
func (p *MusicPlayer) GetCurrentMusicInfo() custom_canvas.MusicInfo {
	return p.musicInfoList[p.musicIndex]
}

/**
Launch me to start the music player!
*/
func (p *MusicPlayer) Launch() {

	volume := 1.0
	for {
		if len(p.musicInfoList) > 0 {

			//load music file
			musicInfo := p.GetCurrentMusicInfo()
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
			p.UpdateProgressBar <- 0.0
			p.UpdateMusicTitle <- musicInfo
			p.UpdatePlayIcon <- true

			player.Play()
			pausedByUser := false
			changeByUser := false
			for {
				if !player.IsPlaying() && !pausedByUser || changeByUser {
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
					p.UpdatePlayIcon <- !pausedByUser

				//roll back song
				case <-p.RequestPrev:
					p.goPrev()
					changeByUser = true

				//skip song
				case <-p.RequestNext:
					p.goNext()
					changeByUser = true

				//reload music list
				case <-p.RequestAdhoc:
					changeByUser = true

				//user set to next playing mode
				case <-p.RequestMode:
					p.nextMode()

				//set to user wanted volume
				case newVolume := <-p.RequestVolume:
					volume = newVolume
					player.SetVolume(volume)

				//set to user wanted music progress
				case percent := <-p.RequestProgress:
					fmt.Print(" ") //magic line to synchronize the program
					player.Pause()
					tick := int64(percent * float64(decodedMp3File.Length()))
					tick -= tick % 4
					if _, err := decodedMp3File.Seek(tick, io.SeekStart); err != nil {
						log.Println(err)
					}
					player.Play()
					p.UpdatePlayIcon <- true

				//update music progress
				default:
					currTick, err := decodedMp3File.Seek(0, io.SeekCurrent)
					if err != nil {
						log.Println(err)
					}
					percent := float64(currTick) / float64(decodedMp3File.Length())
					p.UpdateProgressBar <- percent
					time.Sleep(100.0 * time.Millisecond)
				}
			}

			//normal playing sequence
			if !changeByUser {
				p.goNext()
			}

			//close the resources
			if err = player.Close(); err != nil {
				log.Panic(err)
			}
			if err = mp3File.Close(); err != nil {
				log.Panic(err)
			}

		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
