package scraper

import (
	"fyne.io/fyne/v2/canvas"
)

type ClipzagResult struct {
	videoID     string
	thumbnail   *canvas.Image
	videoTitle  string
	stats       string
	description string
}

func (clipzagResult *ClipzagResult) VideoID() string {
	return clipzagResult.videoID
}

func (clipzagResult *ClipzagResult) Thumbnail() *canvas.Image {
	return clipzagResult.thumbnail
}

func (clipzagResult *ClipzagResult) VideoTitle() string {
	return clipzagResult.videoTitle
}

func (clipzagResult *ClipzagResult) Stats() string {
	return clipzagResult.stats
}

func (clipzagResult *ClipzagResult) Description() string {
	return clipzagResult.description
}
