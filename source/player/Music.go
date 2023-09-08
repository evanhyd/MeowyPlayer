package player

import (
	"fmt"
	"time"
)

const (
	MAGIC_RATIO     = 11024576435 //pray it doesn't overflow
	AUDIO_BIT_DEPTH = 2
	NUM_OF_CHANNELS = 2
	SAMPLING_RATE   = 44100
)

type Music struct {
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	FileSize int64     `json:"-"`
}

func (m *Music) Length() time.Duration {
	return time.Duration(m.FileSize * MAGIC_RATIO / (AUDIO_BIT_DEPTH * NUM_OF_CHANNELS * SAMPLING_RATE))
}

func (m *Music) Description() string {
	const (
		kConversionFactor = 60
		kExtensionLength  = 4 // "remove .mp3"
	)

	length := m.Length()
	mins := int(length.Minutes()) % kConversionFactor
	secs := int(length.Seconds()) % kConversionFactor
	title := m.Title[:len(m.Title)-kExtensionLength]
	return fmt.Sprintf("%02v:%02v | %v", mins, secs, title)
}
