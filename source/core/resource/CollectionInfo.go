package resource

import "time"

type CollectionInfo struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Size  int64     `json:"size"`
}
