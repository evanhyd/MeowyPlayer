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
	"golang.org/x/exp/slices"
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

var mp3Player MP3Player

func init() {
	context, ready, err := oto.NewContext(SAMPLING_RATE, NUM_OF_CHANNELS, AUDIO_BIT_DEPTH)
	if err != nil {
		log.Fatal(err)
	}
	<-ready

	mp3Player = MP3Player{}
	mp3Player.Context = context
	mp3Player.loadMusicChan = make(chan struct{}, 16)
	mp3Player.playPauseMusicChan = make(chan struct{}, 16)
	mp3Player.musicVolumeChan = make(chan float64, 16)
	mp3Player.progressChan = make(chan float64, 16)
	mp3Player.musicVolume = 1.0
	mp3Player.playMode = RANDOM
}

func GetPlayer() *MP3Player {
	return &mp3Player
}

type MP3Player struct {
	*oto.Context

	isLoaded              bool
	loadMusicChan         chan Signal
	playPauseMusicChan    chan Signal
	musicVolumeChan       chan float64
	progressChan          chan float64
	onMusicBeginSubject   pattern.OneArgObservable[Music]
	onMusicPlayingSubject pattern.TwoArgObservable[Music, float64]

	musics      []Music
	musicIndex  int
	musicVolume float64
	playMode    int

	//random mode only
	randomHistory []int
	randomQueue   []int
}

func (p *MP3Player) OnMusicBeginSubject() pattern.OneArgObservabler[Music] {
	return &p.onMusicBeginSubject
}

func (p *MP3Player) OnMusicPlayingSubject() pattern.TwoArgObservabler[Music, float64] {
	return &p.onMusicPlayingSubject
}

func (p *MP3Player) SetMusic(album Album, musics []Music, music Music) {
	p.randomHistory = []int{}
	p.randomQueue = []int{}
	p.musics = musics

	if index := slices.Index(musics, music); index == -1 {
		log.Fatal("Can not find the music from the album")
	} else {
		p.musicIndex = index
	}

	if !p.isLoaded {
		p.isLoaded = true
	} else {
		p.loadMusicChan <- Signal{}
	}
}

func (p *MP3Player) SetPlayMode(playMode int) {
	if playMode == RANDOM {
		p.randomHistory = []int{}
		p.randomQueue = rand.Perm(len(p.musics))
	}
	p.playMode = playMode
}

func (p *MP3Player) SetProgress(percent float64) {
	if p.isLoaded {
		p.progressChan <- percent
	}
}

func (p *MP3Player) SetMusicVolume(volume float64) {
	p.musicVolume = volume
	if p.isLoaded {
		p.musicVolumeChan <- volume
	}
}

func (p *MP3Player) PlayPauseMusic() {
	if p.isLoaded {
		p.playPauseMusicChan <- Signal{}
	}
}

func (p *MP3Player) PreviousMusic() {
	if p.isLoaded {
		switch p.playMode {
		case RANDOM:
			if len(p.randomHistory) > 0 {
				p.randomQueue = append(p.randomQueue, p.musicIndex)
				lastIndex := len(p.randomHistory) - 1
				p.musicIndex = p.randomHistory[lastIndex]
				p.randomHistory = p.randomHistory[:lastIndex]
			}

		case ORDERED:
			p.musicIndex = (p.musicIndex - 1 + len(p.musics)) % len(p.musics)

		case REPEAT:

		default:
			log.Fatal("Invalid music play mode")
		}

		p.loadMusicChan <- Signal{}
	}
}

func (p *MP3Player) NextMusic() {
	if p.isLoaded {
		switch p.playMode {
		case RANDOM:
			p.randomHistory = append(p.randomHistory, p.musicIndex)
			if len(p.randomQueue) == 0 {
				p.randomQueue = rand.Perm(len(p.musics))
			}
			lastIndex := len(p.randomQueue) - 1
			p.musicIndex = p.randomQueue[lastIndex]
			p.randomQueue = p.randomQueue[:lastIndex]

		case ORDERED:
			p.musicIndex = (p.musicIndex + 1) % len(p.musics)

		case REPEAT:

		default:
			log.Fatal("Invalid music play mode")
		}

		p.loadMusicChan <- Signal{}
	}
}

func (p *MP3Player) Launch() {
	for {
		if p.isLoaded {

			//read music file
			mp3File, err := os.ReadFile(resource.GetMusicPath(p.musics[p.musicIndex].title))
			if err != nil {
				log.Fatal(err)
			}

			//decode music file
			mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3File))
			if err != nil {
				log.Fatal(err)
			}

			//obtain music player
			mp3Player := p.Context.NewPlayer(mp3Decoder)

			//PLAYYYYYYYYYYYYYYYYYYYYYYYYYYYY
			paused := false
			interrupted := false
			p.onMusicBeginSubject.NotifyAll(p.musics[p.musicIndex])
			mp3Player.SetVolume(p.musicVolume)
			mp3Player.Play()

		MusicLoop:
			for mp3Player.IsPlaying() || paused {
				select {
				case <-p.loadMusicChan:
					interrupted = true
					break MusicLoop

				case <-p.playPauseMusicChan:
					if mp3Player.IsPlaying() {
						paused = true
						mp3Player.Pause()
					} else {
						paused = false
						mp3Player.Play()
					}

				case percent := <-p.progressChan:
					mp3Player.Pause()
					tick := int64(float64(mp3Decoder.Length()) * percent)
					tick -= tick % 4
					time.Sleep(10 * time.Millisecond) //mp3 library has race condition bug
					mp3Decoder.Seek(tick, io.SeekStart)
					mp3Player.Play()

				case volume := <-p.musicVolumeChan:
					mp3Player.SetVolume(volume)

				default:
					currentTick, err := mp3Decoder.Seek(0, io.SeekCurrent)
					if err != nil {
						log.Fatal(err)
					}
					percent := float64(currentTick) / float64(mp3Decoder.Length())
					p.onMusicPlayingSubject.NotifyAll(p.musics[p.musicIndex], percent)
					time.Sleep(200 * time.Millisecond)
				}
			}

			if !interrupted {
				p.NextMusic()
			}

			mp3Player.Close()

		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
