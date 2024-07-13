package player

import (
	"io"
	"math/rand"
	"playground/model"
	"playground/pattern"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

const (
	kSamplingRate              = 44100
	kChannelCount              = 2
	kMaxHistorySize            = 128
	kCommandBufferSize         = 128
	kMediaProgressUpdatePeriod = 1000 * time.Millisecond
)

type MP3PlayerMode int
type MP3PlayerCommand = func()

const (
	KRandomMode MP3PlayerMode = iota
	KOrderedMode
	KRepeatMode
)

type MP3Player struct {
	music        []model.Music
	playIndex    int
	history      []model.Music
	historyIndex int
	queue        []int
	queueIndex   int
	mode         MP3PlayerMode
	commands     chan MP3PlayerCommand

	context *oto.Context
	decoder *mp3.Decoder
	player  *oto.Player

	onAlbumPlayed     pattern.Subject[model.AlbumKey]
	onMusicPlayed     pattern.Subject[model.Music]
	onProgressUpdated pattern.Subject[float64]
}

var mp3Player MP3Player

func InitPlayer() {
	context, ready, err := oto.NewContext(&oto.NewContextOptions{SampleRate: kSamplingRate, ChannelCount: kChannelCount, Format: oto.FormatSignedInt16LE})
	if err != nil {
		fyne.LogError("failed to initialize mp3 context", err)
	}
	<-ready

	mp3Player = MP3Player{
		commands:          make(chan MP3PlayerCommand, kCommandBufferSize),
		context:           context,
		onAlbumPlayed:     pattern.MakeSubject[model.AlbumKey](),
		onMusicPlayed:     pattern.MakeSubject[model.Music](),
		onProgressUpdated: pattern.MakeSubject[float64](),
	}
	mp3Player.run()
}

func Instance() *MP3Player                                          { return &mp3Player }
func (p *MP3Player) OnAlbumPlayed() pattern.Subject[model.AlbumKey] { return p.onAlbumPlayed }
func (p *MP3Player) OnMusicPlayed() pattern.Subject[model.Music]    { return p.onMusicPlayed }
func (p *MP3Player) OnProgressUpdated() pattern.Subject[float64]    { return p.onProgressUpdated }

func (p *MP3Player) Play() {
	p.commands <- func() {
		if p.player.IsPlaying() {
			p.player.Pause()
		} else {
			p.player.Play()
		}
	}
}

func (p *MP3Player) Prev() {
	p.commands <- func() {
		p.historyIndex = max(0, p.historyIndex-1)
		music := p.history[p.historyIndex]
		p.loadMusic(music)
	}
}

func (p *MP3Player) SetMode(mode MP3PlayerMode) {
	p.mode = mode
}

func (p *MP3Player) playNext() {
	//read from history queue
	p.historyIndex = min(len(p.history), p.historyIndex+1)
	if p.historyIndex < len(p.history) {
		music := p.history[p.historyIndex]
		p.loadMusic(music)
		return
	}

	//read from play queue
	switch p.mode {
	case KRandomMode:
		p.playIndex = p.queue[p.queueIndex]
		p.queueIndex = (p.queueIndex + 1) % len(p.queue)

	case KOrderedMode:
		p.playIndex = (p.playIndex + 1) % len(p.music)

	case KRepeatMode:
	}
	music := p.music[p.playIndex]
	p.appendToHistory(music)
	p.loadMusic(music)
}

func (p *MP3Player) Next() {
	p.commands <- p.playNext
}

func (p *MP3Player) SetProgress(percent float64) {
	p.commands <- func() {
		bytes := int64(float64(p.decoder.Length()) * percent)
		bytes -= bytes % 4

		//mp3 library has race conditions, this is a quick fix
		p.player.Pause()
		time.Sleep(10 * time.Millisecond)
		if _, err := p.player.Seek(bytes, io.SeekStart); err != nil {
			fyne.LogError("mp3 set progress failed", err)
		}
		time.Sleep(10 * time.Millisecond)
		p.player.Play()
	}
}

func (p *MP3Player) SetVolume(percent float64) {
	p.commands <- func() {
		p.player.SetVolume(percent * percent)
	}
}

func (p *MP3Player) LoadAlbum(key model.AlbumKey, music []model.Music, index int) {
	p.music = music
	p.playIndex = index
	p.historyIndex = len(p.history) //move out of the history queue
	p.shuffleQueue(index)           //reset random play queue
	p.playNext()
	p.onAlbumPlayed.NotifyAll(key)
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
	p.decoder, err = mp3.NewDecoder(reader)
	if err != nil {
		return err
	}

	//create mp3 player
	if p.player != nil {
		p.player.Close()
	}
	p.player = p.context.NewPlayer(p.decoder)
	p.player.Play()
	p.onMusicPlayed.NotifyAll(music)
	return nil
}

func (p *MP3Player) shuffleQueue(firstToPlay int) {
	p.queue, p.queueIndex = rand.Perm(len(p.music)), 0
	toPlay := slices.Index(p.queue, firstToPlay)
	p.queue[0], p.queue[toPlay] = p.queue[toPlay], p.queue[0]
}

func (p *MP3Player) progress() float64 {
	//must use decoder.Seek() instead of player
	//player.Seek() clears the internal buffer and stutters
	pos, err := p.decoder.Seek(0, io.SeekCurrent)
	if err != nil {
		fyne.LogError("failed to get progress", err)
	}
	return float64(pos) / float64(p.decoder.Length())
}

func (p *MP3Player) isOver() bool {
	pos, err := p.decoder.Seek(0, io.SeekCurrent)
	if err != nil {
		fyne.LogError("failed to get progress", err)
	}
	return pos == p.decoder.Length()
}

func (p *MP3Player) run() {
	go func() {
		updateTimer := time.NewTicker(kMediaProgressUpdatePeriod)
		for {
			select {
			case cmd := <-p.commands:
				if p.player != nil {
					cmd()
				}
			case <-updateTimer.C:
				if p.player != nil {
					p.onProgressUpdated.NotifyAll(p.progress())
					if p.isOver() {
						p.Next()
					}
				}
			}
		}
	}()
}
