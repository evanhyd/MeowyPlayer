package players

import (
	"log"
	"log/slog"
	"math/rand"
	"meowyplayer/events"
	"meowyplayer/storages"
	"slices"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

const (
	targetSampleRate beep.SampleRate = 48000
)

var _ MusicPlayer = &BeepPlayer{}

type BeepPlayer struct {
	sync.Mutex

	historyStack  []storages.Music
	playlistQueue []storages.Music
	playlistIndex int
	isPlaying     bool
	isRepeating   bool
	queueMode     QueueMode

	storage             storages.Storage
	stream              *BeepStream
	finishedPlayingChan chan struct{}
}

func MakeBeepPlayer() *BeepPlayer {
	go sync.OnceFunc(func() {
		// Set the buffer size delay to be lower than human perception time.
		err := speaker.Init(targetSampleRate, targetSampleRate.N(100*time.Millisecond))
		if err != nil {
			log.Fatalf("failed to initialize the beep speaker: %v\n", err)
		}
	})()

	beepPlayer := &BeepPlayer{
		playlistIndex:       0,
		isPlaying:           false,
		isRepeating:         false,
		queueMode:           SequentialQueueMode,
		finishedPlayingChan: make(chan struct{}),
	}
	go beepPlayer.finishedPlayingRoutine()
	return beepPlayer
}

func (p *BeepPlayer) Resume() {
	p.Lock()
	defer p.Unlock()
	if p.stream != nil {
		p.stream.resume()
	}
}

func (p *BeepPlayer) Suspend() {
	p.Lock()
	defer p.Unlock()
	if p.stream != nil {
		p.stream.suspend()
	}
}

func (p *BeepPlayer) Previous() {
	p.Lock()
	defer p.Unlock()

	// Do nothing if empty playlist.
	if len(p.playlistQueue) == 0 {
		return
	}

	// Reach the end of the historyStack, do nothing.
	if -p.playlistIndex == len(p.historyStack) {
		return
	}
	p.playlistIndex--

	p.playCurrentMusic()
}

func (p *BeepPlayer) Next() {
	p.Lock()
	defer p.Unlock()

	// Do nothing if empty playlist.
	if len(p.playlistQueue) == 0 {
		return
	}

	// Add to the play history.
	if len(p.historyStack) >= 128 {
		p.historyStack = p.historyStack[64:]
	}
	p.historyStack = append(p.historyStack, p.playlistQueue[p.playlistIndex])

	// Wrap around if reach the end of the playlist.
	p.playlistIndex++
	if p.playlistIndex == len(p.playlistQueue) {
		p.playlistIndex = 0
	}

	p.playCurrentMusic()
}

func (p *BeepPlayer) SetQueueMode(queueMode QueueMode) {
	p.Lock()
	defer p.Unlock()
	p.queueMode = queueMode
	if p.queueMode == RandomQueueMode {
		p.shufflePlaylist()
	}
}

func (p *BeepPlayer) SetRepeat(isRepeating bool) {
	p.Lock()
	defer p.Unlock()
	p.isRepeating = isRepeating
}

func (p *BeepPlayer) SetProgress(progress float64) {
	p.Lock()
	defer p.Unlock()
	if p.stream != nil {
		p.stream.setProgress(progress)
	}
}

func (p *BeepPlayer) SetVolume(volume float64) {
	p.Lock()
	defer p.Unlock()
	if p.stream != nil {
		p.stream.setVolume(volume)
	}
}

func (p *BeepPlayer) getCurrentMusic() storages.Music {
	// negative -> historyStack, positive -> playlistQueue
	if p.playlistIndex >= 0 {
		return p.playlistQueue[p.playlistIndex]
	} else {
		return p.historyStack[-p.playlistIndex-1]
	}
}

func (p *BeepPlayer) shufflePlaylist() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(p.playlistQueue), func(i, j int) { p.playlistQueue[i], p.playlistQueue[j] = p.playlistQueue[j], p.playlistQueue[i] })
}

func (p *BeepPlayer) playCurrentMusic() {
	// Get music through the storage.
	content, err := p.storage.GetMusicFile(p.getCurrentMusic())
	if err != nil {
		slog.Error("failed to get music file from the storage", "error", err)
		return
	}

	// Create a stream and play.
	p.stream = newBeepStream(content, targetSampleRate)
	speaker.Clear()
	speaker.Play(beep.Seq(p.stream, beep.Callback(func() {
		p.finishedPlayingChan <- struct{}{}
	})))
}

func (p *BeepPlayer) finishedPlayingRoutine() {
	for range p.finishedPlayingChan {
		if p.isRepeating {
			p.Lock()
			p.playCurrentMusic()
			p.Unlock()
		} else {
			p.Next()
		}
	}
}

// Event Handler
func (p *BeepPlayer) HandleStorageSetEvent(event events.EventType, eventData any) {
	p.Lock()
	defer p.Unlock()
	p.storage = eventData.(events.StorageSetEventData).Storage
}

func (p *BeepPlayer) HandlePlaylistSetEvent(event events.EventType, eventData any) {
	p.Lock()
	defer p.Unlock()
	data := eventData.(events.PlaylistSetEventData)
	p.historyStack = p.historyStack[:0]
	p.playlistQueue = data.MusicList
	if p.queueMode == RandomQueueMode {
		p.shufflePlaylist()
	}
	p.playlistIndex = slices.Index(p.playlistQueue, data.MusicList[data.Index])
	p.playCurrentMusic()
}
