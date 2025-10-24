package player

import (
	"meowyplayer/storage"
	"sync"
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

	stream *BeepStream
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

}

func (p *BeepPlayer) Next() {

}

func (p *BeepPlayer) SetQueueMode(QueueMode) {

}

func (p *BeepPlayer) SetRepeat(isRepeating bool) {

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

func (p *BeepPlayer) SetPlaylist([]storage.Music) {

}
