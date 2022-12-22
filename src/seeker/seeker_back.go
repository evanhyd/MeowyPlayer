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

	RequestRollBack chan struct{}
	RequestSkip     chan struct{}
	RequestAdhoc    chan struct{}
	RequestPlay     chan struct{}
	RequestMode     chan struct{}
	RequestVolume   chan float64
	RequestProgress chan float64

	otoCtx        *oto.Context
	musicInfoList []custom_canvas.MusicInfo
	musicIndex    int
	mode          int   //music player mode: normal, loop, random
	history       []int //played music history
	randomQueue   []int //random music queue
}

func NewMusicPlayer() (*MusicPlayer, error) {
	p := &MusicPlayer{}
	p.UpdateMusicTitle = make(chan custom_canvas.MusicInfo, 4)
	p.UpdateProgressBar = make(chan float64, 4)
	p.UpdatePlayIcon = make(chan bool, 4)
	p.UpdateModeIcon = make(chan int, 4)

	p.RequestRollBack = make(chan struct{}, 4)
	p.RequestSkip = make(chan struct{}, 4)
	p.RequestAdhoc = make(chan struct{}, 4)
	p.RequestPlay = make(chan struct{}, 4)
	p.RequestMode = make(chan struct{}, 4)
	p.RequestVolume = make(chan float64, 4)
	p.RequestProgress = make(chan float64, 4)

	var ready chan struct{}
	var err error
	p.otoCtx, ready, err = oto.NewContext(SAMPLING_RATE, NUM_OF_CHANNELS, AUDIO_BIT_DEPTH)
	if err != nil {
		return nil, err
	}
	<-ready

	p.musicInfoList = make([]custom_canvas.MusicInfo, 0)
	p.musicIndex = 0
	p.mode = RANDOM_MODE
	p.history = make([]int, 0)
	p.randomQueue = make([]int, 0)
	return p, nil
}

/**
Set the music queue to the provided music info.
Reload and play the music at the given index.
*/
func (p *MusicPlayer) SetPlaylist(musicInfo []custom_canvas.MusicInfo, musicIndex int) {
	p.musicInfoList = musicInfo
	p.musicIndex = musicIndex
	p.history = make([]int, 0)
	p.randomQueue = make([]int, 0)
	p.RequestAdhoc <- struct{}{}
}

/**
Roll back to the previous song.
If the mode is "loop", then it rolls back to the previous song in the music playlist.
If the mode is "repeat",
	if roll back by the user, then it rolls back to the previous song in the music playlist.
  if the song naturally ended, then it repeats
If the mode is "random", then it rolls back to the previous song in the played history.
*/
func (p *MusicPlayer) rollBack(roll_back_by_user bool) {
	switch p.mode {
	case LOOP_MODE:
		p.musicIndex = (p.musicIndex - 1 + len(p.musicInfoList)) % len(p.musicInfoList)

	case REPEAT_MODE:
		if roll_back_by_user {
			p.musicIndex = (p.musicIndex - 1 + len(p.musicInfoList)) % len(p.musicInfoList)
		}

	case RANDOM_MODE:
		if len(p.history) > 0 {
			p.musicIndex, p.history = p.history[len(p.history)-1], p.history[:len(p.history)-1]
		}
	}
}

/**
Skip to the next song.
If the mode is "loop", then it skips to the next song in the music playlist.
If the mode is "repeat",
	if skip by the user, then it skips to the next song in the music playlist.
  if the song naturally ended, then it repeats
If the mode is "random", then it skips to a random song that's predetermined in the random music queue.
*/
func (p *MusicPlayer) skip(skip_by_user bool) {
	switch p.mode {
	case LOOP_MODE:
		p.musicIndex = (p.musicIndex + 1) % len(p.musicInfoList)

	case REPEAT_MODE:
		if skip_by_user {
			p.musicIndex = (p.musicIndex + 1) % len(p.musicInfoList)
		}

	case RANDOM_MODE:
		p.history = append(p.history, p.musicIndex)
		if len(p.randomQueue) == 0 {
			p.randomQueue = rand.Perm(len(p.musicInfoList))
		}
		p.musicIndex, p.randomQueue = p.randomQueue[0], p.randomQueue[1:]
	}
}

/**
Change the playing mode to the next mode.
*/
func (p *MusicPlayer) nextMode() {
	p.mode = (p.mode + 1) % MODE_LEN
	p.UpdateModeIcon <- p.mode
}

/**
Get the current music info
*/
func (p *MusicPlayer) getCurrentMusicInfo() custom_canvas.MusicInfo {
	return p.musicInfoList[p.musicIndex]
}

/**
Main thread to handle the music player logic.
*/
func (p *MusicPlayer) Launch() {

	user_volume := 1.0
	for {
		if len(p.musicInfoList) > 0 {

			playMusic := func() {
				//load music file
				musicInfo := p.getCurrentMusicInfo()
				mp3File, err := os.Open(resource.GetMusicPath(musicInfo.Title))
				if err != nil {
					log.Panic(err)
				}
				defer mp3File.Close()

				//decode music file
				decodedMp3File, err := mp3.NewDecoder(mp3File)
				if err != nil {
					log.Panic(err)
				}

				//obtain music player
				player := p.otoCtx.NewPlayer(decodedMp3File)
				defer player.Close()

				//updaet music volume, GUI
				player.SetVolume(user_volume)
				p.UpdateMusicTitle <- musicInfo
				p.UpdatePlayIcon <- true
				pausedByUser := false

				for player.Play(); player.IsPlaying() || pausedByUser; {
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
					case <-p.RequestRollBack:
						p.rollBack(true)
						return

					//skip song
					case <-p.RequestSkip:
						p.skip(true)
						return

					//reload music list
					case <-p.RequestAdhoc:
						return

					//user set to next playing mode
					case <-p.RequestMode:
						p.nextMode()

					//set to user wanted volume
					case user_volume = <-p.RequestVolume:
						player.SetVolume(user_volume)

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

				//go to next song naturally based on the player mode
				p.skip(false)
			}
			playMusic()

		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
