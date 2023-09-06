package player

import (
	"time"
)

type Config struct {
	Date   time.Time `json:"date"`
	Albums []Album   `json:"albums"`
}
