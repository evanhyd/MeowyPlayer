package player

import (
	"fmt"
	"io"
	"meowyplayer/model"
	"meowyplayer/scraper"
	"meowyplayer/util"
	"time"

	"fyne.io/fyne/v2"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

const kSampleRate beep.SampleRate = 48000

type MP3PlayerCommand = func()

type Mp3Stream struct {
	resource      beep.StreamSeekCloser
	resamplerCtrl *beep.Resampler
	volumeCtrl    *effects.Volume
	playCtrl      *beep.Ctrl
}

func newMp3Stream(rsc io.ReadSeekCloser, percent float64) (*Mp3Stream, error) {
	resource, format, err := mp3.Decode(rsc)
	if err != nil {
		return nil, err
	}
	resamplerCtrl := beep.Resample(10, format.SampleRate, kSampleRate, resource)
	volumeCtrl := &effects.Volume{Streamer: resamplerCtrl, Base: 10, Volume: 2 * (percent - 0.8)}
	playCtrl := &beep.Ctrl{Streamer: volumeCtrl}
	return &Mp3Stream{resource, resamplerCtrl, volumeCtrl, playCtrl}, nil
}

func (s *Mp3Stream) Stream(sample [][2]float64) (int, bool) {
	return s.playCtrl.Stream(sample)
}

func (s *Mp3Stream) Close() error {
	return s.resource.Close()
}

func (s *Mp3Stream) Err() error {
	return s.playCtrl.Err()
}

type MP3Player struct {
	queue         MusicQueue
	commands      chan MP3PlayerCommand
	stream        *Mp3Stream
	volumePercent float64

	onAlbumPlayed     util.Subject[model.AlbumKey]
	onMusicPlayed     util.Subject[model.Music]
	onProgressUpdated util.Subject[float64]
}

var mp3Player MP3Player

func InitPlayer() error {
	const kCommandBufferSize = 128
	const kDefaultVolume = 0.5

	//human perception time
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
		speaker.Clear() //remove the callback to avoid double close
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
		p.volumePercent = percent
		speaker.Lock()
		p.stream.volumeCtrl.Volume = 2 * (percent - 0.8) //default to Base^0 == 1 scaling
		p.stream.volumeCtrl.Silent = (percent == 0.0)
		speaker.Unlock()
	}
}

func (p *MP3Player) SetProgress(percent float64) {
	p.commands <- func() {
		speaker.Lock()
		bytes := int(float64(p.stream.resource.Len()) * percent)
		bytes -= bytes % 4
		p.stream.resource.Seek(bytes)
		p.stream.playCtrl.Paused = false
		speaker.Unlock()
	}
}

func (p *MP3Player) LoadAlbum(key model.AlbumKey, musicQueue []model.Music, index int) {
	speaker.Clear()
	p.loadMusic(p.queue.loadPlaylist(musicQueue, index))
	p.onAlbumPlayed.NotifyAll(key)
}

func (p *MP3Player) loadMusic(music *model.Music) {
	reader, err := model.StorageClient().GetMusic(music.Key())
	if err != nil {
		badMusicLogInfo := fmt.Sprintf("%s, %s", music.Title(), music.Key())
		fyne.LogError("synchronizing missing music: "+badMusicLogInfo, err)

		//sync music content to the storage
		if err := model.StorageClient().SyncMusic(scraper.Result{Platform: music.Platform(), ID: music.ID()}); err != nil {
			fyne.LogError("failed to sync music: "+badMusicLogInfo, err)
			return
		}

		//reload
		reader, err = model.StorageClient().GetMusic(music.Key())
		if err != nil {
			fyne.LogError("detected missing music after sync: "+badMusicLogInfo, err)
			return
		}
	}

	//create mp3 stream and controllers
	p.stream, err = newMp3Stream(reader, p.volumePercent)
	if err != nil {
		badMusicLogInfo := fmt.Sprintf("%s, %s", music.Title(), music.Key())
		fyne.LogError("detected bad mp3 file: "+badMusicLogInfo, err)
		return
	}
	speaker.Play(beep.Seq(p.stream, beep.Callback(func() {
		p.stream.Close()
		p.Next()
	})))
	p.onMusicPlayed.NotifyAll(*music)
}

func (p *MP3Player) getProgressPercent() float64 {
	return float64(p.stream.resource.Position()) / float64(p.stream.resource.Len())
}

func (p *MP3Player) run() {
	go func() {
		const kMediaProgressUpdatePeriod = 1000 * time.Millisecond
		updateTimer := time.NewTicker(kMediaProgressUpdatePeriod)
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
