package client

import (
	"bytes"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/network/fileformat"
)

var musicLock sync.Mutex

func addMusic(music resource.Music, musicData []byte) error {
	musicLock.Lock()
	defer musicLock.Unlock()

	//write data to the music repo
	if err := os.WriteFile(resource.MusicPath(&music), musicData, 0777); err != nil {
		return err
	}

	//add the music info to the album
	album := getSourceAlbum(albumData.Get())
	if !slices.ContainsFunc(album.MusicList, func(m resource.Music) bool { return m.Title == music.Title }) {
		album.MusicList.PushBack(music)
	}
	if err := reloadCollectionData(); err != nil {
		return err
	}
	return reloadAlbumData()
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
	return addMusic(music, musicData)
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

	data, err := os.ReadFile(musicInfo.URI().Path())
	if err != nil {
		return err
	}
	length, err := estimateMP3DataLength(data)
	if err != nil {
		return err
	}
	music := resource.Music{Date: time.Now(), Title: musicInfo.URI().Name(), Length: length}
	return addMusic(music, data)
}

func DeleteMusic(music *resource.Music) error {
	musicLock.Lock()
	defer musicLock.Unlock()
	album := getSourceAlbum(albumData.Get())

	//remove from the collection, but not delete it from the music repo
	index := slices.IndexFunc(album.MusicList, func(m resource.Music) bool { return m.SimpleTitle() == music.SimpleTitle() })
	assert.Ensure(func() bool { return index != -1 })
	album.MusicList.Remove(index)

	if err := reloadCollectionData(); err != nil {
		return err
	}
	return reloadAlbumData()
}
