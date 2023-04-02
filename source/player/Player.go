package player

import (
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
type PlayMode int

const (
	MAGIC_RATIO     = 11024576435 //pray it doesn't overflow
	AUDIO_BIT_DEPTH = 2
	NUM_OF_CHANNELS = 2
	SAMPLING_RATE   = 44100
)

const (
	RANDOM PlayMode = iota
	ORDER
	LOOP
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
	selectMusicUpdater

	isLoaded  bool
	loadMusic chan Signal

	album       Album
	musics      []Music
	musicIndex  int
	musicVolume float64

	playMode      PlayMode
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
	player.selectMusicUpdater = selectMusicUpdater{player}
	player.loadMusic = make(chan struct{}, 16)
	player.musicVolume = 1.0
	player.playMode = RANDOM
	return player
}

func (player *Player) PlayerMusicUpdater() pattern.ThreeArgObserver[Album, []Music, Music] {
	return &player.selectMusicUpdater
}

func (player *Player) ChangePlayMode(playMode PlayMode) {
	if playMode == RANDOM {
		player.randomHistory = []int{}
		player.randomQueue = rand.Perm(player.album.musicNumber)
	}
	player.playMode = playMode
}

func (player *Player) PreviousMusic() {
	switch player.playMode {
	case RANDOM:
		if len(player.randomHistory) > 0 {
			player.randomQueue = append(player.randomQueue, player.musicIndex)
			lastIndex := len(player.randomHistory) - 1
			player.musicIndex = player.randomHistory[lastIndex]
			player.randomHistory = player.randomHistory[:lastIndex]
		}

	case ORDER:
		player.musicIndex = (player.musicIndex - 1 + player.album.musicNumber) % player.album.musicNumber

	case LOOP:

	default:
		log.Fatal("Invalid music play mode")
	}

	player.loadMusic <- Signal{}
}

func (player *Player) NextMusic() {
	switch player.playMode {
	case RANDOM:
		player.randomHistory = append(player.randomHistory, player.musicIndex)
		if len(player.randomQueue) == 0 {
			player.randomQueue = rand.Perm(player.album.musicNumber)
		}
		lastIndex := len(player.randomQueue) - 1
		player.musicIndex = player.randomQueue[lastIndex]
		player.randomQueue = player.randomQueue[:lastIndex]

	case ORDER:
		player.musicIndex = (player.musicIndex + 1) % player.album.musicNumber

	case LOOP:

	default:
		log.Fatal("Invalid music play mode")
	}

	player.loadMusic <- Signal{}
}

func (player *Player) Launch() {
	for {
		if player.isLoaded {

			//open music file
			mp3File, err := os.Open(resource.GetMusicPath(player.musics[player.musicIndex].title))
			if err != nil {
				log.Fatal(err)
			}

			//decode music stream
			mp3Stream, err := mp3.NewDecoder(mp3File)
			if err != nil {
				log.Fatal(err)
			}

			//obtain music player
			mp3Player := player.Context.NewPlayer(mp3Stream)

			//load current config
			mp3Player.SetVolume(player.musicVolume)

			interrupted := false
		MusicLoop:
			for mp3Player.Play(); mp3Player.IsPlaying(); {
				select {
				case <-player.loadMusic:
					interrupted = true
					break MusicLoop
				default:
					log.Println("idling")
					time.Sleep(1000 * time.Millisecond)
				}
			}

			if !interrupted {
				player.NextMusic()
			}

			mp3Player.Close()
			mp3File.Close()

		} else {
			log.Println("waiting to load")
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

type selectMusicUpdater struct {
	*Player
}

func (player *selectMusicUpdater) Notify(album Album, musics []Music, music Music) {
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
	log.Printf("[%v] %v\n", album.title, music.title)

	if !player.isLoaded {
		player.isLoaded = true
	} else {
		player.loadMusic <- Signal{}
	}
}
