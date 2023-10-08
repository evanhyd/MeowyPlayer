package resource

import (
	"time"

	"meowyplayer.com/utility/container"
)

type Collection struct {
	Date   time.Time              `json:"date"`
	Albums container.Slice[Album] `json:"albums"`
}
