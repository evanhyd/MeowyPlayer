package player

import (
	"fmt"
	"meowyplayer/model"
	"meowyplayer/scraper"
	"meowyplayer/util"
	"time"

	"fyne.io/fyne/v2"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

const kSampleRate beep.SampleRate = 48000

type MP3PlayerCommand = func()

type MP3Player struct {
	queue         MusicQueue
	commands      chan MP3PlayerCommand
	stream        *ControllableMp3Stream
	volumePercent float64

	onAlbumPlayed     util.Subject[model.AlbumKey]
	onMusicPlayed     util.Subject[model.Music]
	onProgressUpdated util.Subject[float64]
}

var mp3Player MP3Player

func InitPlayer() error {
	const kCommandBufferSize = 128
	const kDefaultVolume = 0.5

	// Set the buffer size delay to be lower than human perception time.
	if err := speaker.Init(kSampleRate, kSampleRate.N(100*time.Millisecond)); err != nil {
		return err
	}

	mp3Player = MP3Player{
		commands:          make(chan MP3PlayerCommand, kCommandBufferSize),
		volumePercent:     kDefaultVolume,
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

func (p *MP3Player) Prev() {
	p.commands <- func() {
		speaker.Clear() // Remove the callback to avoid double close.
		p.loadMusic(p.queue.prev())
	}
}

func (p *MP3Player) Next() {
	p.commands <- func() {
		speaker.Clear()
		p.loadMusic(p.queue.next())
	}
}

func (p *MP3Player) Play() {
	p.commands <- func() {
		speaker.Lock()
		p.stream.playCtrl.Paused = !p.stream.playCtrl.Paused
		speaker.Unlock()
	}
}

func (p *MP3Player) SetVolume(percent float64) {
	p.commands <- func() {
		speaker.Lock()
		p.volumePercent = percent
		p.stream.setVolume(percent)
		speaker.Unlock()
	}
}

func (p *MP3Player) SetProgress(percent float64) {
	p.commands <- func() {
		speaker.Lock()
		bytes := int(float64(p.stream.resource.Len()) * percent) // Align to 4 bytes.
		bytes -= bytes % 4
		p.stream.resource.Seek(bytes)
		p.stream.playCtrl.Paused = false
		speaker.Unlock()
	}
}

func (p *MP3Player) LoadAlbum(key model.AlbumKey, musicQueue []model.Music, toPlayIndex int) {
	speaker.Clear()
	p.loadMusic(p.queue.loadPlaylist(musicQueue, toPlayIndex))
	p.onAlbumPlayed.NotifyAll(key)
}

func (p *MP3Player) loadMusic(music model.Music) {
	musicLogInfo := fmt.Sprintf("%s, %s", music.Title(), music.Key())
	reader, err := model.StorageClient().GetMusic(music.Key())
	if err != nil {
		fyne.LogError("missing music in the storage: "+musicLogInfo, err)

		// Sync music content to the storage.
		if err := model.StorageClient().SyncMusic(scraper.Result{Platform: music.Platform(), ID: music.ID(), Title: music.Title()}); err != nil {
			fyne.LogError("failed to sync music: "+musicLogInfo, err)
			return
		}

		// Reload
		reader, err = model.StorageClient().GetMusic(music.Key())
		if err != nil {
			fyne.LogError("missing music after sync: "+musicLogInfo, err)
			return
		}
	}

	//create mp3 stream and controllers
	if p.stream, err = newMp3Stream(reader); err != nil {
		fyne.LogError("detected bad mp3 file: "+musicLogInfo, err)
		return
	}

	p.stream.setVolume(p.volumePercent)
	speaker.Play(beep.Seq(p.stream, beep.Callback(func() {
		p.stream.Close()
		p.Next()
	})))
	p.onMusicPlayed.NotifyAll(music)
}

func (p *MP3Player) getProgressPercent() float64 {
	return float64(p.stream.resource.Position()) / float64(p.stream.resource.Len())
}

func (p *MP3Player) run() {
	go func() {
		updateTimer := time.NewTicker(1000 * time.Millisecond)
		for {
			select {
			case cmd := <-p.commands:
				if p.stream != nil {
					cmd()
				}
			case <-updateTimer.C:
				if p.stream != nil {
					p.onProgressUpdated.NotifyAll(p.getProgressPercent())
				}
			}
		}
	}()
}
