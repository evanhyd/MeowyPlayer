package players

import (
	"log/slog"
	"math/rand/v2"
	"meowyplayer/storages"
	"slices"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

const sampleRate beep.SampleRate = 48000

var _ MusicPlayer = (*BeepPlayer)(nil)

type BeepPlayer struct {
	userContext    *storages.UserContext
	onMissingMusic func(storages.Music)
	onPlayingMusic func(storages.Music)

	musicList       ringBuffer[storages.Music]
	randomIndexList ringBuffer[int]
	randomHistory   int
	isRepeating     bool
	queueMode       QueueMode
	volume          float64

	currentMusic storages.Music
	stream       *BeepStream

	actionChan chan func()
}

func MakeBeepPlayer(userContext *storages.UserContext, onMissingMusic, onPlayingMusic func(storages.Music)) *BeepPlayer {
	sync.OnceFunc(func() {
		if err := speaker.Init(sampleRate, sampleRate.N(100*time.Millisecond)); err != nil {
			slog.Error("failed to initialize the beep speaker", "error", err)
		}
	})()

	p := &BeepPlayer{
		userContext:    userContext,
		onMissingMusic: onMissingMusic,
		onPlayingMusic: onPlayingMusic,
		volume:         0.7,
		actionChan:     make(chan func(), 16),
	}

	go func() {
		for action := range p.actionChan {
			action()
		}
	}()

	return p
}

func (p *BeepPlayer) IsPlaying() bool {
	res := make(chan bool, 1)
	p.actionChan <- func() { res <- p.stream != nil && !p.stream.playCtrl.Paused }
	return <-res
}

func (p *BeepPlayer) GetProgress() float64 {
	res := make(chan float64, 1)
	p.actionChan <- func() {
		if p.stream == nil {
			res <- 0.0
		} else {
			res <- p.stream.progress()
		}
	}
	return <-res
}

func (p *BeepPlayer) GetMusic() storages.Music {
	res := make(chan storages.Music, 1)
	p.actionChan <- func() { res <- p.currentMusic }
	return <-res
}

func (p *BeepPlayer) Resume() {
	p.actionChan <- func() {
		if p.stream != nil {
			p.stream.resume()
		}
	}
}

func (p *BeepPlayer) Pause() {
	p.actionChan <- func() {
		if p.stream != nil {
			p.stream.suspend()
		}
	}
}

func (p *BeepPlayer) Previous() {
	p.actionChan <- func() {
		if p.musicList.empty() {
			return
		}
		switch p.queueMode {
		case SequentialQueueMode:
			p.playMusic(p.musicList.prev())
		case RandomQueueMode:
			if p.randomHistory > 0 {
				p.randomHistory--
				p.playMusic(p.musicList.data[p.randomIndexList.prev()])
			}
		}
	}
}

func (p *BeepPlayer) Next() {
	p.actionChan <- func() {
		if p.musicList.empty() {
			return
		}
		if p.isRepeating && p.stream != nil {
			p.playMusic(p.currentMusic)
			return
		}
		switch p.queueMode {
		case SequentialQueueMode:
			p.playMusic(p.musicList.next())
		case RandomQueueMode:
			p.randomHistory++
			p.playMusic(p.musicList.data[p.randomIndexList.next()])
		}
	}
}

func (p *BeepPlayer) SetPlaylist(playlist storages.Playlist, selectedMusic storages.Music) {
	p.actionChan <- func() {
		music, err := p.userContext.GetMusicFromPlaylist(playlist.PlaylistId)
		if err != nil {
			slog.Error("failed to fetch music", "error", err)
			return
		}
		idx := slices.Index(music, selectedMusic)
		if idx == -1 {
			slog.Error("selected music not in playlist", "playlist", playlist.PlaylistId)
			return
		}
		p.musicList.set(music)
		p.musicList.index = idx
		p.queueMode = SequentialQueueMode
		p.playMusic(selectedMusic)
	}
}

func (p *BeepPlayer) SetQueueMode(mode QueueMode) {
	p.actionChan <- func() {
		p.queueMode = mode
		if mode == SequentialQueueMode && len(p.musicList.data) > 0 {
			p.musicList.index = p.randomIndexList.current()
		} else if mode == RandomQueueMode {
			p.randomIndexList.set(rand.Perm(len(p.musicList.data)))
			p.randomHistory = 0
		}
	}
}

func (p *BeepPlayer) SetRepeat(isRepeating bool) {
	p.actionChan <- func() { p.isRepeating = isRepeating }
}

func (p *BeepPlayer) SetProgress(progress float64) {
	p.actionChan <- func() {
		if p.stream != nil {
			p.stream.setProgress(progress)
		}
	}
}

func (p *BeepPlayer) SetVolume(volume float64) {
	p.actionChan <- func() {
		p.volume = volume
		if p.stream != nil {
			p.stream.setVolume(volume)
		}
	}
}

func (p *BeepPlayer) playMusic(music storages.Music) {
	p.currentMusic = music
	musicFile, err := p.userContext.GetMusicFile(music)

	if err != nil {
		go func(m storages.Music) {
			p.onMissingMusic(m)
			p.actionChan <- func() {
				if p.currentMusic.MusicId == m.MusicId {
					p.playMusic(m)
				}
			}
		}(music)
		return
	}

	speaker.Clear()
	if p.stream != nil {
		p.stream.Close()
	}

	p.stream = newBeepStream(musicFile, sampleRate, p.volume)
	if p.stream == nil {
		go p.onMissingMusic(music)
		return
	}

	speaker.Play(beep.Seq(p.stream, beep.Callback(p.Next)))
	go p.onPlayingMusic(music)
}
