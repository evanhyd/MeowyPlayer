package resource

import (
	"fmt"
	"time"
)

const (
	SAMPLING_RATE   = 44100
	NUM_OF_CHANNELS = 2
	AUDIO_BIT_DEPTH = 2
)

type Music struct {
	Date     time.Time     `json:"date"`
	Title    string        `json:"title"`
	Length   time.Duration `json:"length"`
	Platform string        `json:"platform"`
	ID       string        `json:"id"`
}

// return title without the extension string
func (m *Music) SimpleTitle() string {
	const kExtensionLength = 4
	return m.Title[:len(m.Title)-kExtensionLength]
}

func (m *Music) Description() string {
	const kConversionFactor = 60
	mins := int(m.Length.Minutes()) % kConversionFactor
	secs := int(m.Length.Seconds()) % kConversionFactor
	return fmt.Sprintf("%02v:%02v | %v", mins, secs, m.SimpleTitle())
}
