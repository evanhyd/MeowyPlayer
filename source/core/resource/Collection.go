package resource

import (
	"time"
)

type Collection struct {
	Date   time.Time        `json:"date"`
	Albums map[string]Album `json:"albums"`
}
