package client

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/logger"
	"meowyplayer.com/utility/network/downloader"
	"meowyplayer.com/utility/network/fileformat"
)

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
		album.Title = fmt.Sprintf("Album %v", i)
		if err := Manager().addAlbum(album); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to add new album")
}

func isMusicFileExist(music *resource.Music) bool {
	_, err := os.Stat(resource.MusicPath(music))
	return err == nil
}

func AddMusicFromURIReader(album resource.Album, musicInfo fyne.URIReadCloser) error {
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
	return Manager().addMusic(album, resource.Music{Date: time.Now(), Title: musicInfo.URI().Name(), Length: length}, musicData)
}

func DownloadMusic(album resource.Album, videoResult *fileformat.VideoResult) error {
	var provider downloader.MusicDownloader
	switch videoResult.Platform {
	case "YouTube":
		provider = downloader.NewY2MateDownloader()
	case "BiliBili":
		logger.Fatal(fmt.Errorf("not implemented"), 0)
	default:
		return nil
	}

	data, err := provider.Download(videoResult)
	if err != nil {
		return err
	}

	music := resource.Music{
		Date:     time.Now(),
		Title:    resource.SanatizeFileName(videoResult.Title) + ".mp3",
		Length:   videoResult.Length,
		Platform: videoResult.Platform,
		ID:       videoResult.VideoID,
	}
	return Manager().addMusic(album, music, data)
}

func CloneMusic(album resource.Album, music resource.Music) error {
	return DownloadMusic(album, &fileformat.VideoResult{
		Title:    music.Title[:len(music.Title)-4],
		Length:   music.Length,
		Platform: music.Platform,
		VideoID:  music.ID,
	})
}

func SyncCollection() int32 {
	var unsynced atomic.Int32
	wg := sync.WaitGroup{}
	for _, album := range Manager().collection.Get().Albums {
		for _, music := range album.MusicList {
			if !isMusicFileExist(&music) {
				wg.Add(1)
				go func(album resource.Album, music resource.Music) {
					defer wg.Done()
					if err := CloneMusic(album, music); err != nil {
						unsynced.Add(1)
						logger.Error(err, 0)
					}
				}(album, music)
			}
		}
	}
	wg.Wait()
	return unsynced.Load()
}
