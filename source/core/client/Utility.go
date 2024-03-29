package client

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
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
	stat, err := os.Stat(resource.MusicPath(music))
	return err == nil && stat.Size() > 0
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

/*
Download music from the internet based on the videoResult, then add to the album.
*/
func AddMusicFromDownloader(album resource.Album, videoResult *fileformat.VideoResult) error {
	music := resource.Music{
		Date:     time.Now(),
		Title:    resource.SanatizeFileName(videoResult.Title) + ".mp3",
		Length:   videoResult.Length,
		Platform: videoResult.Platform,
		ID:       videoResult.VideoID,
	}

	data, err := downloadMusic(videoResult)
	if err != nil {
		return err
	}

	return Manager().addMusic(album, music, data)
}

func downloadMusic(videoResult *fileformat.VideoResult) ([]byte, error) {
	var provider downloader.MusicDownloader
	switch videoResult.Platform {
	case "YouTube":
		provider = downloader.NewY2MateDownloader()
	case "BiliBili":
		return nil, fmt.Errorf("BiliBili downloader is not implemented")
	default:
		return nil, fmt.Errorf("unknown downloader")
	}

	return provider.Download(videoResult)
}

func SyncCollection() <-chan float64 {
	const kDownloadRoutines = 4
	percents := make(chan float64, kDownloadRoutines)
	toDownload := []resource.Music{}
	for _, album := range Manager().currentCollection.Get().Albums {
		for _, music := range album.MusicList {
			if !isMusicFileExist(&music) {
				toDownload = append(toDownload, music)
			}
		}
	}

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(toDownload))
		tokens := make(chan struct{}, kDownloadRoutines)
		success := atomic.Int32{}
		for _, music := range toDownload {
			go func(music resource.Music) {
				tokens <- struct{}{}
				if err := syncMusic(music); err != nil {
					logger.Error(err, 1)
				} else {
					percents <- float64(success.Add(1)) / float64(len(toDownload))
				}
				<-tokens
				wg.Done()
			}(music)
		}
		wg.Wait()
		close(percents)
		Manager().load()
		log.Printf("%v/%v music downloaded\n", success.Load(), len(toDownload))
	}()
	return percents
}

/*
Download the missing music file.
*/
func syncMusic(music resource.Music) error {
	videoResult := fileformat.VideoResult{
		Title:    music.Title[:len(music.Title)-4],
		Length:   music.Length,
		Platform: music.Platform,
		VideoID:  music.ID,
	}
	data, err := downloadMusic(&videoResult)
	if err != nil {
		return err
	}

	return os.WriteFile(resource.MusicPath(&music), data, 0777)
}
