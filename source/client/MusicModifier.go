package client

import (
	"bytes"
	"os"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/utility/assert"
)

func estimateMP3DataLength(data []byte) time.Duration {
	decoder, err := mp3.NewDecoder(bytes.NewReader(data))
	assert.NoErr(err, "failed to decode mp3 data")
	seconds := float64(decoder.Length()) / float64(resource.SAMPLING_RATE) / float64(resource.NUM_OF_CHANNELS) / float64(resource.AUDIO_BIT_DEPTH)
	return time.Duration(seconds * float64(time.Second))
}

func AddLocalMusic(musicInfo fyne.URIReadCloser) error {
	music := resource.Music{Date: time.Now(), Title: musicInfo.URI().Name()}

	//copy the music file to the music repo
	data, err := os.ReadFile(musicInfo.URI().Path())
	if err != nil {
		return err
	}
	if err = os.WriteFile(resource.MusicPath(&music), data, os.ModePerm); err != nil {
		return err
	}
	music.Length = estimateMP3DataLength(data)

	album := getSourceAlbum(albumData.Get())
	album.MusicList = append(album.MusicList, music)
	if err := reloadCollectionData(); err != nil {
		return err
	}
	return reloadAlbumData()
}

func DeleteMusic(music *resource.Music) error {
	album := getSourceAlbum(albumData.Get())
	index := slices.IndexFunc(album.MusicList, func(m resource.Music) bool { return m.SimpleTitle() == music.SimpleTitle() })
	last := len(album.MusicList) - 1

	//pop form the album
	album.MusicList[index] = album.MusicList[last]
	album.MusicList = album.MusicList[:last]

	if err := reloadCollectionData(); err != nil {
		return err
	}
	return reloadAlbumData()
}
