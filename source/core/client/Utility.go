package client

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/json"
	"meowyplayer.com/utility/network/fileformat"
)

func LoadFromLocalCollection() (resource.Collection, error) {
	collection := resource.Collection{}
	if err := json.ReadFile(resource.CollectionPath(), &collection); err != nil {
		return collection, err
	}

	for title, album := range collection.Albums {
		album.Cover = resource.GetCover(&album)
		collection.Albums[title] = album
	}

	return collection, nil
}

func AddRandomAlbum() error {
	//generate album cover
	iconColor := color.NRGBA{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())}
	iconImage := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	iconImage.SetNRGBA(0, 0, iconColor)
	imageData := bytes.Buffer{}
	if err := png.Encode(&imageData, iconImage); err != nil {
		return err
	}

	//generate album
	album := resource.Album{
		Date:      time.Now(),
		Title:     "",
		MusicList: make(map[string]resource.Music),
		Cover:     fyne.NewStaticResource("", imageData.Bytes()),
	}

	// try 100 possible titles until it fits in
	for i := 0; i < 100; i++ {
		album.Title = fmt.Sprintf("Album (%v)", i)
		if err := GetInstance().AddAlbum(album); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to add new album")
}

func AddMusicFromDownloader(videoResult *fileformat.VideoResult, musicData []byte) error {
	//sanitize music title
	sanitizer := strings.NewReplacer(
		"<", "",
		">", "",
		":", "",
		"\"", "",
		"/", "",
		"\\", "",
		"|", "",
		"?", "",
		"*", "",
	)
	music := resource.Music{Date: time.Now(), Title: sanitizer.Replace(videoResult.Title) + ".mp3", Length: videoResult.Length}
	return GetInstance().AddMusic(music, musicData)
}

func AddMusicFromURIReader(musicInfo fyne.URIReadCloser) error {
	estimateMP3DataLength := func(data []byte) (time.Duration, error) {
		decoder, err := mp3.NewDecoder(bytes.NewReader(data))
		if err != nil {
			return 0, err
		}
		seconds := float64(decoder.Length()) / float64(resource.SAMPLING_RATE) / float64(resource.NUM_OF_CHANNELS) / float64(resource.AUDIO_BIT_DEPTH)
		return time.Duration(seconds * float64(time.Second)), nil
	}

	musicData, err := os.ReadFile(musicInfo.URI().Path())
	if err != nil {
		return err
	}
	length, err := estimateMP3DataLength(musicData)
	if err != nil {
		return err
	}
	music := resource.Music{Date: time.Now(), Title: musicInfo.URI().Name(), Length: length}
	return GetInstance().AddMusic(music, musicData)
}
