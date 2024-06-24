package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var _ FileSystem = &localStorage{}

type localStorage struct {
	sync.RWMutex
	albumDir string
	musicDir string
}

func NewLocalStorage() *localStorage {
	const kStorage = "storage"
	return &localStorage{
		albumDir: filepath.Join(kStorage, "album"),
		musicDir: filepath.Join(kStorage, "music"),
	}
}

func (f *localStorage) albumPath(key AlbumKey) string {
	return filepath.Join(f.albumDir, fmt.Sprintf("%v.json", key))
}

func (f *localStorage) musicPath(key MusicKey) string {
	return filepath.Join(f.musicDir, fmt.Sprintf("%v.mp3", key))
}

func (f *localStorage) initialize() error {
	if err := os.MkdirAll(f.albumDir, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(f.musicDir, 0700); err != nil {
		return err
	}
	return nil
}

func (f *localStorage) getAllAlbums() ([]Album, error) {
	f.RLock()
	defer f.RUnlock()

	const kFileExt = ".json"
	entries, err := os.ReadDir(f.albumDir)
	if err != nil {
		return nil, err
	}

	albums := make([]Album, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == kFileExt {
			data, err := os.ReadFile(filepath.Join(f.albumDir, entry.Name()))
			if err != nil {
				return nil, err
			}

			var album Album
			if err := json.Unmarshal(data, &album); err != nil {
				return nil, err
			}
			albums = append(albums, album)
		}
	}
	return albums, nil
}

func (f *localStorage) getAlbum(key AlbumKey) (album Album, err error) {
	if key.IsEmpty() {
		return album, fmt.Errorf("empty key in getAlbum")
	}

	f.RLock()
	defer f.RUnlock()

	data, err := os.ReadFile(f.albumPath(key))
	if err == nil {
		err = json.Unmarshal(data, &album)
	}
	return
}

func (f *localStorage) uploadAlbum(album Album) error {
	key := album.Key()
	if key.IsEmpty() {
		return fmt.Errorf("empty key in uploadAlbum")
	}

	f.Lock()
	defer f.Unlock()

	data, err := json.Marshal(&album)
	if err != nil {
		return err
	}
	return os.WriteFile(f.albumPath(key), data, 0600)
}

func (f *localStorage) removeAlbum(key AlbumKey) error {
	if key.IsEmpty() {
		return fmt.Errorf("empty key in removeAlbum")
	}

	f.Lock()
	defer f.Unlock()

	return os.Remove(f.albumPath(key))
}

func (f *localStorage) getMusic(key MusicKey) (io.ReadSeekCloser, error) {
	if key.IsEmpty() {
		return nil, fmt.Errorf("empty key in getMusic")
	}

	f.RLock()
	defer f.RUnlock()

	return os.Open(f.musicPath(key))
}

func (f *localStorage) uploadMusic(music Music, reader io.Reader) error {
	key := music.Key()
	if key.IsEmpty() {
		return fmt.Errorf("empty key in uploadMusic")
	}

	f.Lock()
	defer f.Unlock()

	dst, err := os.OpenFile(f.musicPath(key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, reader)
	return err
}

func (f *localStorage) removeMusic(key MusicKey) error {
	if key.IsEmpty() {
		return fmt.Errorf("empty key in removeMusic")
	}

	f.Lock()
	defer f.Unlock()

	return os.Remove(f.musicPath(key))
}

// func (d *localFileSystem) sanatizeFileName(filename string) string {
// 	sanitizer := strings.NewReplacer(
// 		"<", "",
// 		">", "",
// 		":", "",
// 		"\"", "",
// 		"/", "",
// 		"\\", "",
// 		"|", "",
// 		"?", "",
// 		"*", "",
// 		"~", "",
// 	)
// 	return sanitizer.Replace(filename)
// }
