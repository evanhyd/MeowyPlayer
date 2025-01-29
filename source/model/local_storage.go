package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
)

var _ Storage = &localStorage{}

type localStorage struct {
	sync.Mutex
	albumDir string
	musicDir string
}

func newLocalStorage() *localStorage {
	const kStorage = "storage"
	s := localStorage{
		albumDir: filepath.Join(kStorage, "local"),
		musicDir: filepath.Join(kStorage, "music"),
	}
	if err := os.MkdirAll(s.albumDir, 0700); err != nil {
		fyne.LogError("can not create local storage album dir", err)
	}
	if err := os.MkdirAll(s.musicDir, 0700); err != nil {
		fyne.LogError("can not create music dir", err)
	}
	return &s
}

func (s *localStorage) albumPath(key AlbumKey) string {
	return filepath.Join(s.albumDir, fmt.Sprintf("%v.json", key))
}

func (s *localStorage) musicPath(key MusicKey) string {
	return filepath.Join(s.musicDir, fmt.Sprintf("%v.mp3", key))
}

func (s *localStorage) getAllAlbums() ([]Album, error) {
	s.Lock()
	defer s.Unlock()

	const kFileExt = ".json"
	entries, err := os.ReadDir(s.albumDir)
	if err != nil {
		return nil, err
	}

	albums := make([]Album, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == kFileExt {
			data, err := os.ReadFile(filepath.Join(s.albumDir, entry.Name()))
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

func (s *localStorage) getAlbum(key AlbumKey) (Album, error) {
	s.Lock()
	defer s.Unlock()

	if key.IsEmpty() {
		return Album{}, fmt.Errorf("empty key in getAlbum")
	}

	data, err := os.ReadFile(s.albumPath(key))
	if err != nil {
		return Album{}, err
	}

	album := Album{}
	err = json.Unmarshal(data, &album)
	return album, err
}

func (s *localStorage) uploadAlbum(album Album) error {
	s.Lock()
	defer s.Unlock()

	key := album.Key()
	if key.IsEmpty() {
		return fmt.Errorf("empty key in uploadAlbum")
	}

	data, err := json.Marshal(album)
	if err != nil {
		return err
	}
	return os.WriteFile(s.albumPath(key), data, 0600)
}

func (s *localStorage) removeAlbum(key AlbumKey) error {
	s.Lock()
	defer s.Unlock()

	if key.IsEmpty() {
		return fmt.Errorf("empty key in removeAlbum")
	}
	return os.Remove(s.albumPath(key))
}

func (s *localStorage) getMusic(key MusicKey) (io.ReadSeekCloser, error) {
	s.Lock()
	defer s.Unlock()

	if key.IsEmpty() {
		return nil, fmt.Errorf("empty key in getMusic")
	}
	return os.Open(s.musicPath(key))
}

func (s *localStorage) uploadMusic(music Music, content io.Reader) error {
	s.Lock()
	defer s.Unlock()

	key := music.Key()
	if key.IsEmpty() {
		return fmt.Errorf("empty key in uploadMusic")
	}

	file, err := os.OpenFile(s.musicPath(key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

func (s *localStorage) removeMusic(key MusicKey) error {
	s.Lock()
	defer s.Unlock()

	if key.IsEmpty() {
		return fmt.Errorf("empty key in removeMusic")
	}
	return os.Remove(s.musicPath(key))
}

// func (d *localStorage) sanatizeFileName(filename string) string {
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
