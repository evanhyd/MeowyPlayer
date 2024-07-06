package player

import (
	"math/rand"
	"playground/model"
	"slices"

	"fyne.io/fyne/v2"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

const (
	kSamplingRate      = 44100
	kChannelCount      = 2
	kMaxHistorySize    = 128
	kCommandBufferSize = 128
)

type MP3PlayerMode int

const (
	kRandomMode MP3PlayerMode = iota
	kOrderedMode
	kRepeatMode
)

type MP3Player struct {
	music        []model.Music
	playIndex    int
	history      []model.Music
	historyIndex int
	queue        []int
	queueIndex   int
	mode         MP3PlayerMode

	commands chan MP3Command
	context  *oto.Context
	player   *oto.Player
	mp3Size  int64
}

var player *MP3Player

func Instance() *MP3Player {
	return player
}

func InitPlayer() {
	context, ready, err := oto.NewContext(&oto.NewContextOptions{SampleRate: kSamplingRate, ChannelCount: kChannelCount, Format: oto.FormatSignedInt16LE})
	if err != nil {
		fyne.LogError("failed to initialize mp3 context", err)
	}
	<-ready

	player = &MP3Player{
		commands: make(chan MP3Command, kCommandBufferSize),
		context:  context,
	}
}

func (p *MP3Player) Play()                        { p.commands <- mp3Play{} }
func (p *MP3Player) Prev()                        { p.commands <- mp3Prev{} }
func (p *MP3Player) Next()                        { p.commands <- mp3Next{} }
func (p *MP3Player) SetProgress(progress float64) { p.commands <- mp3Progress{progress} }
func (p *MP3Player) SetVolume(volume float64)     { p.commands <- m3pVolume{volume} }

func (p *MP3Player) LoadAlbum(music []model.Music, index int) {
	p.music = music
	p.playIndex = index
	p.historyIndex = len(p.history) //move out of the history queue
	p.shuffleQueue(index)           //reset random play queue
	mp3Next{}.execute(p)
}

func (p *MP3Player) appendToHistory(music model.Music) {
	if len(p.history) == 0 || music != p.history[len(p.history)-1] {
		if len(p.history) >= kMaxHistorySize {
			p.history = p.history[1:]
		}
		p.history = append(p.history, music)
		p.historyIndex = len(p.history)
	}
}

func (p *MP3Player) loadMusic(music model.Music) error {
	//decode mp3
	reader, err := model.Instance().GetMusic(music.Key())
	if err != nil {
		return err
	}

	decoder, err := mp3.NewDecoder(reader)
	if err != nil {
		return err
	}

	if p.player != nil {
		p.player.Close()
	}
	p.player = p.context.NewPlayer(decoder)
	p.mp3Size = decoder.Length()
	p.player.Play()
	return nil
}

func (p *MP3Player) shuffleQueue(firstToPlay int) {
	p.queue, p.queueIndex = rand.Perm(len(p.music)), 0
	toPlay := slices.Index(p.queue, firstToPlay)
	p.queue[0], p.queue[toPlay] = p.queue[toPlay], p.queue[0]
}

func (p *MP3Player) run() {
	for cmd := range p.commands {
		if p.player != nil {
			cmd.execute(p)
		}
	}
}
