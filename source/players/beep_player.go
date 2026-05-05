package players

import (
	"log/slog"
	"math/rand/v2"
	"meowyplayer/context"
	"meowyplayer/storages"
	"slices"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

const (
	sampleRate beep.SampleRate = 48000
)

// ==========================================
// Command Interface
// ==========================================
type Command interface{}
type CmdIsPlaying struct{ isPlaying chan bool }
type CmdGetProgress struct{ progress chan float64 }
type CmdGetMusic struct{ music chan storages.Music }
type CmdResume struct{}
type CmdPause struct{}
type CmdPlayNext struct{}
type CmdPlayPrev struct{}
type CmdSetQueueMode struct{ mode QueueMode }
type CmdSetRepeat struct{ isRepeating bool }
type CmdSetProgress struct{ percent float64 }
type CmdSetVolume struct{ volume float64 }
type CmdSetPlaylist struct {
	playlist storages.Playlist
	music    storages.Music
}

// ==========================================
// BeepPlayer Implementation
// ==========================================

var _ MusicPlayer = &BeepPlayer{}

type BeepPlayer struct {
	userContext     *context.UserContext
	musicList       ringBuffer[storages.Music]
	randomIndexList ringBuffer[int]
	randomHistory   int
	isRepeating     bool
	queueMode       QueueMode
	volume          float64

	currentMusic storages.Music
	stream       *BeepStream

	cmdChan chan Command
}

func MakeBeepPlayer(userContext *context.UserContext) *BeepPlayer {
	sync.OnceFunc(func() {
		err := speaker.Init(sampleRate, sampleRate.N(100*time.Millisecond))
		if err != nil {
			slog.Error("failed to initialize the beep speaker", "error", err)
		}
	})()

	p := BeepPlayer{
		userContext: userContext,
		volume:      0.7,
		cmdChan:     make(chan Command, 16),
	}
	go p.run()
	return &p
}

// ---------------------------------------------------------
// Public API: Non-blocking command dispatch
// ---------------------------------------------------------
func (p *BeepPlayer) IsPlaying() bool {
	isPlaying := make(chan bool)
	p.cmdChan <- CmdIsPlaying{isPlaying: isPlaying}
	return <-isPlaying
}

func (p *BeepPlayer) GetProgress() float64 {
	progress := make(chan float64)
	p.cmdChan <- CmdGetProgress{progress: progress}
	return <-progress
}

func (p *BeepPlayer) GetMusic() storages.Music {
	music := make(chan storages.Music)
	p.cmdChan <- CmdGetMusic{music: music}
	return <-music
}

func (p *BeepPlayer) Resume() {
	p.cmdChan <- CmdResume{}
}

func (p *BeepPlayer) Pause() {
	p.cmdChan <- CmdPause{}
}

func (p *BeepPlayer) Previous() {
	p.cmdChan <- CmdPlayPrev{}
}

func (p *BeepPlayer) Next() {
	p.cmdChan <- CmdPlayNext{}
}

func (p *BeepPlayer) SetPlaylist(playlist storages.Playlist, selectedMusic storages.Music) {
	p.cmdChan <- CmdSetPlaylist{playlist: playlist, music: selectedMusic}
}

func (p *BeepPlayer) SetQueueMode(queueMode QueueMode) {
	p.cmdChan <- CmdSetQueueMode{mode: queueMode}
}

func (p *BeepPlayer) SetRepeat(isRepeating bool) {
	p.cmdChan <- CmdSetRepeat{isRepeating: isRepeating}
}

func (p *BeepPlayer) SetProgress(progress float64) {
	p.cmdChan <- CmdSetProgress{percent: progress}
}

func (p *BeepPlayer) SetVolume(volume float64) {
	p.cmdChan <- CmdSetVolume{volume: volume}
}

// ---------------------------------------------------------
// Internal Actor Loop
// ---------------------------------------------------------

func (p *BeepPlayer) playMusic(music storages.Music) {
	p.currentMusic = music
	musicFile, err := p.userContext.Storage().GetMusicFile(music)
	if err != nil {
		slog.Error("failed to get music file from the storage", "error", err)
		return
	}

	speaker.Clear()
	if p.stream != nil {
		p.stream.Close()
	}

	p.stream = newBeepStream(musicFile, sampleRate, p.volume)
	if p.stream == nil {
		return
	}

	speaker.Play(beep.Seq(p.stream, beep.Callback(p.Next)))
}

func (p *BeepPlayer) run() {
	for cmd := range p.cmdChan {
		switch c := cmd.(type) {

		case CmdIsPlaying:
			c.isPlaying <- (p.stream != nil && !p.stream.playCtrl.Paused)

		case CmdGetProgress:
			if p.stream == nil {
				c.progress <- 0.0
			} else {
				c.progress <- p.stream.progress()
			}

		case CmdGetMusic:
			c.music <- p.currentMusic

		case CmdResume:
			if p.stream != nil {
				p.stream.resume()
			}

		case CmdPause:
			if p.stream != nil {
				p.stream.suspend()
			}

		case CmdPlayPrev:
			if p.musicList.empty() {
				continue
			}
			switch p.queueMode {
			case SequentialQueueMode:
				p.playMusic(p.musicList.prev())
			case RandomQueueMode:
				if p.randomHistory > 0 {
					p.randomHistory -= 1
					p.playMusic(p.musicList.data[p.randomIndexList.prev()])
				}
			}

		case CmdPlayNext:
			if p.isRepeating && p.stream != nil {
				p.playMusic(p.currentMusic)
				continue
			}

			if p.musicList.empty() {
				continue
			}
			switch p.queueMode {
			case SequentialQueueMode:
				p.playMusic(p.musicList.next())
			case RandomQueueMode:
				p.randomHistory += 1
				p.playMusic(p.musicList.data[p.randomIndexList.next()])
			}

		case CmdSetPlaylist:
			music, err := p.userContext.Storage().GetAllSortedMusicFromPlaylist(c.playlist.PlaylistId)
			if err != nil {
				slog.Error("failed to fetch music", "error", err)
				continue
			}
			p.musicList.set(music)
			p.musicList.index = slices.Index(p.musicList.data, c.music)
			if p.musicList.index == -1 {
				slog.Error("the selected music is not part of the playlist", "playlist", c.playlist.PlaylistId)
				continue
			}
			p.queueMode = SequentialQueueMode
			p.playMusic(c.music)

		case CmdSetQueueMode:
			p.queueMode = c.mode
			switch p.queueMode {
			case SequentialQueueMode:
				if len(p.musicList.data) > 0 {
					p.musicList.index = p.randomIndexList.current()
				}
			case RandomQueueMode:
				p.randomIndexList.set(rand.Perm(len(p.musicList.data)))
				p.randomHistory = 0
			}

		case CmdSetRepeat:
			p.isRepeating = c.isRepeating

		case CmdSetProgress:
			if p.stream != nil {
				p.stream.setProgress(c.percent)
			}

		case CmdSetVolume:
			p.volume = c.volume
			if p.stream != nil {
				p.stream.setVolume(c.volume)
			}
		}
	}
}
