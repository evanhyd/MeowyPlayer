package player

import (
	"io"
	"log"
	"math/rand"
	"meowyplayer/storage"
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

	historyStack  []storage.Music
	playlistQueue []storage.Music
	playlistIndex int

	isPlaying   bool
	isRepeating bool
	queueMode   QueueMode

	contentGetter       func(storage.Music) io.ReadSeekCloser
	stream              *BeepStream
	finishedPlayingChan chan struct{}
}

func makeBeepPlayer(contentGetter func(storage.Music) io.ReadSeekCloser) *BeepPlayer {
	go sync.OnceFunc(func() {
		// Set the buffer size delay to be lower than human perception time.
		err := speaker.Init(targetSampleRate, targetSampleRate.N(100*time.Millisecond))
		if err != nil {
			log.Fatalf("failed to initialize the beep speaker: %v\n", err)
		}
	})()

	beepPlayer := &BeepPlayer{
		playlistIndex: 0,
		isPlaying:     false,
		isRepeating:   false,
		queueMode:     SequentialQueueMode,

		contentGetter:       contentGetter,
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

func (p *BeepPlayer) getCurrentMusic() storage.Music {
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
	p.stream = newBeepStream(p.contentGetter(p.getCurrentMusic()), targetSampleRate)
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
func (p *BeepPlayer) OnSelectPlaylist(musicList []storage.Music, index int) {
	p.Lock()
	defer p.Unlock()
	p.historyStack = p.historyStack[:0]
	p.playlistQueue = musicList
	if p.queueMode == RandomQueueMode {
		p.shufflePlaylist()
	}
	p.playlistIndex = slices.Index(p.playlistQueue, musicList[index])
	p.playCurrentMusic()
}
