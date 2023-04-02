package player

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/resource"
)

type Signal = struct{}

const (
	MAGIC_RATIO     = 11024576435 //pray it doesn't overflow
	AUDIO_BIT_DEPTH = 2
	NUM_OF_CHANNELS = 2
	SAMPLING_RATE   = 44100
)

const (
	RANDOM = iota
	ORDERED
	REPEAT
	PLAYMODE_LEN
)

var player *Player

func init() {
	player = NewPlayer()
}

func GetPlayer() *Player {
	return player
}

type Player struct {
	*oto.Context

	isLoaded           bool
	loadMusicChan      chan Signal
	playPauseMusicChan chan Signal
	musicVolumeChan    chan float64
	progressChan       chan float64
	onMusicBegin       pattern.OneArgSubject[Music]
	onMusicPlaying     pattern.TwoArgSubject[Music, float64]

	album       Album
	musics      []Music
	musicIndex  int
	musicVolume float64
	playMode    int

	//random mode only
	randomHistory []int
	randomQueue   []int
}

func NewPlayer() *Player {
	context, ready, err := oto.NewContext(SAMPLING_RATE, NUM_OF_CHANNELS, AUDIO_BIT_DEPTH)
	if err != nil {
		log.Fatal(err)
	}
	<-ready

	player := &Player{}
	player.Context = context
	player.loadMusicChan = make(chan struct{}, 16)
	player.playPauseMusicChan = make(chan struct{}, 16)
	player.musicVolumeChan = make(chan float64, 16)
	player.progressChan = make(chan float64, 16)

	player.musicVolume = 1.0
	player.playMode = RANDOM
	return player
}

func (player *Player) OnMusicBegin() *pattern.OneArgSubject[Music] {
	return &player.onMusicBegin
}

func (player *Player) OnMusicPlaying() *pattern.TwoArgSubject[Music, float64] {
	return &player.onMusicPlaying
}

func (player *Player) SetMusic(album Album, musics []Music, music Music) {
	if player.album != album {
		player.album = album
		player.randomHistory = []int{}
		player.randomQueue = []int{}
	}
	player.musics = musics //musics sorting order can be different

	found := false
	for i := range musics {
		if musics[i] == music {
			player.musicIndex = i
			found = true
			break
		}
	}
	if !found {
		log.Fatal("Can not find the music from the album")
	}

	if !player.isLoaded {
		player.isLoaded = true
	} else {
		player.loadMusicChan <- Signal{}
	}
}

func (player *Player) SetPlayMode(playMode int) {
	if playMode == RANDOM {
		player.randomHistory = []int{}
		player.randomQueue = rand.Perm(player.album.musicNumber)
	}
	player.playMode = playMode
}

func (player *Player) SetProgress(percent float64) {
	if player.isLoaded {
		player.progressChan <- percent
	}
}

func (player *Player) SetMusicVolume(volume float64) {
	player.musicVolume = volume
	if player.isLoaded {
		player.musicVolumeChan <- volume
	}
}

func (player *Player) PlayPauseMusic() {
	if player.isLoaded {
		player.playPauseMusicChan <- Signal{}
	}
}

func (player *Player) PreviousMusic() {
	if player.isLoaded {
		switch player.playMode {
		case RANDOM:
			if len(player.randomHistory) > 0 {
				player.randomQueue = append(player.randomQueue, player.musicIndex)
				lastIndex := len(player.randomHistory) - 1
				player.musicIndex = player.randomHistory[lastIndex]
				player.randomHistory = player.randomHistory[:lastIndex]
			}

		case ORDERED:
			player.musicIndex = (player.musicIndex - 1 + player.album.musicNumber) % player.album.musicNumber

		case REPEAT:

		default:
			log.Fatal("Invalid music play mode")
		}

		player.loadMusicChan <- Signal{}
	}
}

func (player *Player) NextMusic() {
	if player.isLoaded {
		switch player.playMode {
		case RANDOM:
			player.randomHistory = append(player.randomHistory, player.musicIndex)
			if len(player.randomQueue) == 0 {
				player.randomQueue = rand.Perm(player.album.musicNumber)
			}
			lastIndex := len(player.randomQueue) - 1
			player.musicIndex = player.randomQueue[lastIndex]
			player.randomQueue = player.randomQueue[:lastIndex]

		case ORDERED:
			player.musicIndex = (player.musicIndex + 1) % player.album.musicNumber

		case REPEAT:

		default:
			log.Fatal("Invalid music play mode")
		}

		player.loadMusicChan <- Signal{}
	}
}

func (player *Player) Launch() {
	for {
		if player.isLoaded {

			//read music file
			mp3File, err := os.ReadFile(resource.GetMusicPath(player.musics[player.musicIndex].title))
			if err != nil {
				log.Fatal(err)
			}

			//decode music file
			mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3File))
			if err != nil {
				log.Fatal(err)
			}

			//obtain music player
			mp3Player := player.Context.NewPlayer(mp3Decoder)

			//PLAYYYYYYYYYYYYYYYYYYYYYYYYYYYY
			paused := false
			interrupted := false
			player.onMusicBegin.NotifyAll(player.musics[player.musicIndex])
			mp3Player.SetVolume(player.musicVolume)
			mp3Player.Play()

		MusicLoop:
			for mp3Player.IsPlaying() || paused {
				select {
				case <-player.loadMusicChan:
					interrupted = true
					break MusicLoop

				case <-player.playPauseMusicChan:
					if mp3Player.IsPlaying() {
						paused = true
						mp3Player.Pause()
					} else {
						paused = false
						mp3Player.Play()
					}

				case percent := <-player.progressChan:
					mp3Player.Pause()
					tick := int64(float64(mp3Decoder.Length()) * percent)
					tick -= tick % 4
					time.Sleep(10 * time.Millisecond) //mp3 library has race condition bug
					mp3Decoder.Seek(tick, io.SeekStart)
					mp3Player.Play()

				case volume := <-player.musicVolumeChan:
					mp3Player.SetVolume(volume)

				default:
					currentTick, err := mp3Decoder.Seek(0, io.SeekCurrent)
					if err != nil {
						log.Fatal(err)
					}
					percent := float64(currentTick) / float64(mp3Decoder.Length())
					player.onMusicPlaying.NotifyAll(player.musics[player.musicIndex], percent)
					time.Sleep(200 * time.Millisecond)
				}
			}

			if !interrupted {
				player.NextMusic()
			}

			mp3Player.Close()

		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
