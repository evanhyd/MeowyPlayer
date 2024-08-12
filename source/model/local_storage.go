package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ Storage = &localStorage{}

type localStorage struct {
	albumDir string
	musicDir string
}

func newLocalStorage() *localStorage {
	const kStorage = "storage"
	return &localStorage{
		albumDir: filepath.Join(kStorage, "album"),
		musicDir: filepath.Join(kStorage, "music"),
	}
}

func (s *localStorage) initialize() error {
	if err := os.MkdirAll(s.albumDir, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(s.musicDir, 0700); err != nil {
		return err
	}
	return nil
}

func (s *localStorage) albumPath(key AlbumKey) string {
	return filepath.Join(s.albumDir, fmt.Sprintf("%v.json", key))
}

func (s *localStorage) musicPath(key MusicKey) string {
	return filepath.Join(s.musicDir, fmt.Sprintf("%v.mp3", key))
}

func (s *localStorage) getAllAlbums() ([]Album, error) {
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

func (s *localStorage) getAlbum(key AlbumKey) (album Album, err error) {
	if key.IsEmpty() {
		return album, fmt.Errorf("empty key in getAlbum")
	}

	data, err := os.ReadFile(s.albumPath(key))
	if err == nil {
		err = json.Unmarshal(data, &album)
	}
	return
}

func (s *localStorage) uploadAlbum(album Album) error {
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
	if key.IsEmpty() {
		return fmt.Errorf("empty key in removeAlbum")
	}

	return os.Remove(s.albumPath(key))
}

func (s *localStorage) getMusic(key MusicKey) (io.ReadSeekCloser, error) {
	if key.IsEmpty() {
		return nil, fmt.Errorf("empty key in getMusic")
	}

	return os.Open(s.musicPath(key))
}

func (s *localStorage) uploadMusic(music Music, reader io.Reader) error {
	key := music.Key()
	if key.IsEmpty() {
		return fmt.Errorf("empty key in uploadMusic")
	}

	dst, err := os.OpenFile(s.musicPath(key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, reader)
	return err
}

func (s *localStorage) removeMusic(key MusicKey) error {
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
