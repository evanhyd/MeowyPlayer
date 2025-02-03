package player

import (
	"math/rand"
	"meowyplayer/model"
	"slices"
)

type QueueMode int

const (
	KRandomMode QueueMode = iota
	KOrderedMode
	KRepeatMode
)

/*
musicQueue queue the music in user-specified order.

historyQueue tracks the played music. The last element in the queue is the most recent music.

randomQueue specifies the index of the music that will get played if the player is in random mode.

If the historyIndex < len(historyQueue), then prev() and next() fetch the music from the historyQueue.
Otherwise, prev() fetch the music from the historyQueue, but next() fetch the music by the queueMode and append it to the historyQueue:

RandomMode: get music from the randomQueue

OrderedMode: choose the next music from the musicQueue depends on the previous musicIndex

RepeatMode: repeat the current music
*/
type MusicQueue struct {
	musicQueue      []model.Music
	lastPlayedIndex int

	historyQueue []model.Music
	historyIndex int

	randomQueue []int
	randomIndex int

	mode QueueMode
}

func (q *MusicQueue) loadPlaylist(musicQueue []model.Music, index int) model.Music {
	q.musicQueue = musicQueue
	q.lastPlayedIndex = index
	q.appendToHistory(index)
	q.shuffleQueue(index)
	return q.musicQueue[q.lastPlayedIndex]
}

func (q *MusicQueue) appendToHistory(index int) {
	// Append if different than the last music in queue.
	if len(q.historyQueue) == 0 || q.musicQueue[index] != q.historyQueue[len(q.historyQueue)-1] {
		q.historyQueue = append(q.historyQueue, q.musicQueue[index])
	}
	q.historyIndex = len(q.historyQueue) - 1
}

func (q *MusicQueue) shuffleQueue(toPlay int) {
	q.randomQueue, q.randomIndex = rand.Perm(len(q.musicQueue)), 0
	index := slices.Index(q.randomQueue, toPlay)
	q.randomQueue[0], q.randomQueue[index] = q.randomQueue[index], q.randomQueue[0]
}

func (q *MusicQueue) setMode(mode QueueMode) {
	q.mode = mode
}

func (q *MusicQueue) prev() model.Music {
	q.historyIndex = max(0, q.historyIndex-1)
	return q.historyQueue[q.historyIndex]
}

func (q *MusicQueue) next() model.Music {
	q.historyIndex = min(len(q.historyQueue), q.historyIndex+1)

	// Browing the history queue.
	if q.historyIndex < len(q.historyQueue) {
		return q.historyQueue[q.historyIndex]
	}

	// Read from the music queue.
	switch q.mode {
	case KRandomMode:
		q.randomIndex = (q.randomIndex + 1) % len(q.randomQueue)
		q.lastPlayedIndex = q.randomQueue[q.randomIndex]
	case KOrderedMode:
		q.lastPlayedIndex = (q.lastPlayedIndex + 1) % len(q.musicQueue)
	case KRepeatMode:
		//nothing
	}
	q.appendToHistory(q.lastPlayedIndex)
	return q.musicQueue[q.lastPlayedIndex]
}
