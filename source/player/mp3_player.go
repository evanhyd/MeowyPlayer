package player

import (
	"io"
	"meowyplayer/browser"
	"meowyplayer/model"
	"meowyplayer/util"
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
	volume  float64

	onAlbumPlayed     util.Subject[model.AlbumKey]
	onMusicPlayed     util.Subject[model.Music]
	onProgressUpdated util.Subject[float64]
}

var mp3Player MP3Player

func InitPlayer() error {
	const (
		kSamplingRate      = 44100
		kChannelCount      = 2
		kCommandBufferSize = 128
		kDefaultVolume     = 0.5
	)

	context, ready, err := oto.NewContext(&oto.NewContextOptions{SampleRate: kSamplingRate, ChannelCount: kChannelCount, Format: oto.FormatSignedInt16LE})
	if err != nil {
		return err
	}
	<-ready

	mp3Player = MP3Player{
		commands:          make(chan MP3PlayerCommand, kCommandBufferSize),
		context:           context,
		volume:            kDefaultVolume,
		onAlbumPlayed:     util.MakeSubject[model.AlbumKey](),
		onMusicPlayed:     util.MakeSubject[model.Music](),
		onProgressUpdated: util.MakeSubject[float64](),
	}
	mp3Player.run()
	return nil
}

func Instance() *MP3Player                                       { return &mp3Player }
func (p *MP3Player) OnAlbumPlayed() util.Subject[model.AlbumKey] { return p.onAlbumPlayed }
func (p *MP3Player) OnMusicPlayed() util.Subject[model.Music]    { return p.onMusicPlayed }
func (p *MP3Player) OnProgressUpdated() util.Subject[float64]    { return p.onProgressUpdated }
func (p *MP3Player) SetMode(mode QueueMode)                      { p.queue.setMode(mode) }
func (p *MP3Player) Prev()                                       { p.commands <- func() { p.loadMusic(p.queue.prev()) } }
func (p *MP3Player) Next()                                       { p.commands <- func() { p.loadMusic(p.queue.next()) } }

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
	p.commands <- func() {
		p.volume = percent * percent
		p.player.SetVolume(p.volume)
	}
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
	reader, err := model.StorageClient().GetMusic(music.Key())
	if err != nil {
		//sync music content to the storage
		if err := model.StorageClient().SyncMusic(browser.Result{Platform: music.Platform(), ID: music.ID()}); err != nil {
			fyne.LogError("failed to sync music", err)
			return
		}

		//reload
		reader, err = model.StorageClient().GetMusic(music.Key())
		if err != nil {
			fyne.LogError("failed to get music content", err)
			return
		}
	}

	p.decoder, err = mp3.NewDecoder(reader)
	if err != nil {
		fyne.LogError("failed to decode music content", err)
		return
	}

	//create mp3 player
	if p.player != nil {
		p.player.Close()
	}
	p.player = p.context.NewPlayer(p.decoder)
	p.player.SetVolume(p.volume)
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
