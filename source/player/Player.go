package player

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"

	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/resource"
)

var state *playerState

func init() {
	state = &playerState{}
}

func GetPlayerState() *playerState {
	return state
}

type playerState struct {
	allAlbums     []Album
	selectedAlbum Album
	accessLock    sync.Mutex

	onUpdateAllAlbums pattern.OneArgSubject[[]Album]
	onSelectAlbum     pattern.OneArgSubject[Album]
}

func (state *playerState) OnUpdateAllAlbumsAddObserver(observer pattern.OneArgObserver[[]Album]) {
	state.accessLock.Lock()
	defer state.accessLock.Unlock()
	state.onUpdateAllAlbums.AddObserver(observer)
}

func (state *playerState) onSelectAlbumAddObserver(observer pattern.OneArgObserver[Album]) {
	state.accessLock.Lock()
	defer state.accessLock.Unlock()
	state.onSelectAlbum.AddObserver(observer)
}

func (state *playerState) UpdateAlbums() {
	state.accessLock.Lock()
	defer state.accessLock.Unlock()

	directories, err := os.ReadDir(resource.GetAlbumPath())
	if err != nil {
		log.Fatal(err)
	}

	const bufferSize = 32 * 1024 //arbitrary magic number
	buffer := make([]byte, bufferSize)
	state.allAlbums = []Album{}

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

			state.allAlbums = append(state.allAlbums, Album{directory.Name(), info.ModTime(), musicNumber, albumCover})
		}
	}

	state.onUpdateAllAlbums.NotifyAll(state.allAlbums)
}
