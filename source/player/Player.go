package player

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/resource"
)

var state *State

func init() {
	state = &State{}
	state.info.state = state
}

func GetState() *State {
	return state
}

type StateInfo struct {
	state *State
	album Album
	music []Music
}

func (info *StateInfo) Notify(album Album) {
	info.album = album
	info.music = ReadMusicFromDirectory(album)
	info.state.onSelectAlbum.NotifyAll(info.album, info.music)
}

type State struct {
	info          StateInfo
	onReadAlbums  pattern.OneArgSubject[[]Album]
	onSelectAlbum pattern.TwoArgSubject[Album, []Music]
}

func (state *State) Info() pattern.OneArgObserver[Album] {
	return &state.info
}

func (state *State) OnReadAlbums() *pattern.OneArgSubject[[]Album] {
	return &state.onReadAlbums
}

func (state *State) OnSelectAlbum() *pattern.TwoArgSubject[Album, []Music] {
	return &state.onSelectAlbum
}

func ReadAlbumsFromDirectory() []Album {
	directories, err := os.ReadDir(resource.GetAlbumFolderPath())
	if err != nil {
		log.Fatal(err)
	}

	const bufferSize = 32 * 1024 //arbitrary magic number
	buffer := make([]byte, bufferSize)
	albums := []Album{}

	for _, directory := range directories {
		if directory.IsDir() {

			//read album config
			config, err := os.Open(resource.GetAlbumConfigPath(directory.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer config.Close()

			//read last modified date
			info, err := config.Stat()
			if err != nil {
				log.Fatal(err)
			}

			//read number of music
			bytesRead, err := config.Read(buffer)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			musicNumber := bytes.Count(buffer[:bytesRead], []byte{'\n'})

			//get album cover
			albumCover := canvas.NewImageFromFile(resource.GetAlbumIconPath(directory.Name()))

			albums = append(albums, Album{directory.Name(), info.ModTime(), musicNumber, albumCover})
		}
	}
	return albums
}

func ReadMusicFromDirectory(album Album) []Music {
	config, err := os.Open(resource.GetAlbumConfigPath(album.title))
	if err != nil {
		log.Fatal(err)
	}
	defer config.Close()

	//read music name from config
	music := []Music{}
	scanner := bufio.NewScanner(config)
	for scanner.Scan() {

		//open music file
		info, err := os.Stat(resource.GetMusicPath(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}

		music = append(music, Music{scanner.Text(), estimateDuration(info), info.ModTime()})
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	return music
}

func estimateDuration(musicFileInfo fs.FileInfo) time.Duration {
	const (
		MAGIC_RATIO     = 11024576435 //pray it doesn't overflow
		AUDIO_BIT_DEPTH = 2
		NUM_OF_CHANNELS = 2
		SAMPLING_RATE   = 44100
	)

	//a very rough estimation of the music duration in nanoseconds
	return time.Duration(musicFileInfo.Size() * MAGIC_RATIO / (AUDIO_BIT_DEPTH * NUM_OF_CHANNELS * SAMPLING_RATE))
}
