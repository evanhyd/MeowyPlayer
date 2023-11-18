package client

import (
	"fmt"
	"os"

	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/logger"
	"meowyplayer.com/utility/network/downloader"
	"meowyplayer.com/utility/network/fileformat"
)

func downloadToLocal(provider downloader.MusicDownloader, music *fileformat.VideoResult) error {
	data, err := provider.Download(music)
	if err != nil {
		return err
	}
	return Manager().AddMusicFromDownloader(music, data)
}

func getDownloader(platform string) downloader.MusicDownloader {
	switch platform {
	case "YouTube":
		return downloader.NewY2MateDownloader()
	case "BiliBili":
		logger.Fatal(fmt.Errorf("not implemented"), 0)
	default:
	}
	return nil
}

func DownloadMusicFromMusic(music *resource.Music) error {
	return downloadToLocal(getDownloader(music.Platform), &fileformat.VideoResult{
		Title:    music.Title,
		Length:   music.Length,
		Platform: music.Platform,
		VideoID:  music.ID,
	})
}

func DownloadMusicFromVideo(video *fileformat.VideoResult) error {
	return downloadToLocal(getDownloader(video.Platform), video)
}

func isMusicExist(music *resource.Music) bool {
	_, err := os.Stat(resource.MusicPath(music))
	return err == nil || os.IsNotExist(err)
}
