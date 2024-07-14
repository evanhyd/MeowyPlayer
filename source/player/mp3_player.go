package player

import (
	"io"
	"playground/model"
	"playground/pattern"
	"time"

	"fyne.io/fyne/v2"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type MP3PlayerCommand = func()

type MP3Player struct {
	queue    MusicQueue
	commands chan MP3PlayerCommand

	context *oto.Context
	decoder *mp3.Decoder
	player  *oto.Player

	onAlbumPlayed     pattern.Subject[model.AlbumKey]
	onMusicPlayed     pattern.Subject[model.Music]
	onProgressUpdated pattern.Subject[float64]
}

var mp3Player MP3Player

func InitPlayer() {
	const (
		kSamplingRate      = 44100
		kChannelCount      = 2
		kCommandBufferSize = 128
	)

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
func (p *MP3Player) SetMode(mode QueueMode)                         { p.queue.setMode(mode) }
func (p *MP3Player) Prev()                                          { p.commands <- func() { p.loadMusic(p.queue.prev()) } }
func (p *MP3Player) Next()                                          { p.commands <- func() { p.loadMusic(p.queue.next()) } }

func (p *MP3Player) Play() {
	p.commands <- func() {
		if p.player.IsPlaying() {
			p.player.Pause()
		} else {
			p.player.Play()
		}
	}
}

func (p *MP3Player) SetVolume(percent float64) {
	p.commands <- func() { p.player.SetVolume(percent * percent) }
}

func (p *MP3Player) SetProgress(percent float64) {
	p.commands <- func() {
		bytes := int64(float64(p.decoder.Length()) * percent)
		bytes -= bytes % 4

		//mp3 library has race conditions
		//delay after pause to sync up the buffer
		p.player.Pause()
		time.Sleep(10 * time.Millisecond)
		if _, err := p.player.Seek(bytes, io.SeekStart); err != nil {
			fyne.LogError("mp3 set progress failed", err)
		}
		time.Sleep(10 * time.Millisecond)
		p.player.Play()
	}
}

func (p *MP3Player) LoadAlbum(key model.AlbumKey, musicQueue []model.Music, index int) {
	p.loadMusic(p.queue.loadPlaylist(musicQueue, index))
	p.onAlbumPlayed.NotifyAll(key)
}

func (p *MP3Player) loadMusic(music *model.Music) {
	//decode mp3
	reader, err := model.Instance().GetMusic(music.Key())
	if err != nil {
		fyne.LogError("failed to get music reader", err)
		return
	}
	p.decoder, err = mp3.NewDecoder(reader)
	if err != nil {
		fyne.LogError("failed to get decode music reader", err)
		return
	}

	//create mp3 player
	if p.player != nil {
		p.player.Close()
	}
	p.player = p.context.NewPlayer(p.decoder)
	p.player.Play()
	p.onMusicPlayed.NotifyAll(*music)
}

func (p *MP3Player) getProgressPercent() float64 {
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
		const kMediaProgressUpdatePeriod = 1000 * time.Millisecond
		updateTimer := time.NewTicker(kMediaProgressUpdatePeriod)
		for {
			select {
			case cmd := <-p.commands:
				if p.player != nil {
					cmd()
				}
			case <-updateTimer.C:
				if p.player != nil {
					p.onProgressUpdated.NotifyAll(p.getProgressPercent())
					if p.isOver() {
						p.Next()
					}
				}
			}
		}
	}()
}
