package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

var _ FileSystem = &localFileSystem{}

type localFileSystem struct {
	albumDir string
	musicDir string

	cache map[AlbumKey]Album
}

func NewLocalFileSystem() localFileSystem {
	const storage = "storage"
	return localFileSystem{
		albumDir: filepath.Join(storage, "album"),
		musicDir: filepath.Join(storage, "music"),
		cache:    map[AlbumKey]Album{},
	}
}

func (f *localFileSystem) initialize() error {
	if err := os.MkdirAll(f.albumDir, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(f.musicDir, 0700); err != nil {
		return err
	}
	return f.load()
}

func (f *localFileSystem) getAllAlbums() ([]Album, error) {
	return maps.Values(f.cache), nil
}

func (f *localFileSystem) getAlbum(key AlbumKey) (Album, error) {
	album, exist := f.cache[key]
	if !exist {
		return album, fmt.Errorf("invalid album key")
	}
	return album, nil
}

func (f *localFileSystem) createAlbum(album Album) (AlbumKey, error) {
	key := AlbumKey(uuid.NewString())
	album.key = key
	album.date = time.Now()
	f.cache[key] = album
	return key, f.save(key)
}

func (f *localFileSystem) updateAlbum(album Album) error {
	_, exist := f.cache[album.key]
	if !exist {
		return fmt.Errorf("invalid album key")
	}
	album.date = time.Now()
	f.cache[album.key] = album
	return f.save(album.key)
}

func (f *localFileSystem) removeAlbum(key AlbumKey) error {
	delete(f.cache, key)
	return f.save(key)
}

func (f *localFileSystem) getMusic(key MusicKey) (io.ReadSeekCloser, error) {
	return nil, nil
}

func (f *localFileSystem) uploadMusic(music Music, src io.Reader) (MusicKey, error) {
	key := music.Key()

	dst, err := os.OpenFile(f.getMusicPath(key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return key, err
}

func (f *localFileSystem) removeMusic(key MusicKey) error {
	return os.Remove(f.getMusicPath(key))
}

func (f *localFileSystem) getAlbumPath(key AlbumKey) string {
	return filepath.Join(f.albumDir, fmt.Sprintf("%v.json", key))
}

func (f *localFileSystem) getMusicPath(key MusicKey) string {
	return filepath.Join(f.musicDir, fmt.Sprintf("%v.mp3", key))
}

func (f *localFileSystem) save(key AlbumKey) error {
	album, exist := f.cache[key]
	if exist {
		file, err := os.OpenFile(f.getAlbumPath(key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		return json.NewEncoder(file).Encode(&album)
	} else {
		return os.Remove(f.getAlbumPath(key))
	}
}

func (f *localFileSystem) load() error {
	entries, err := os.ReadDir(f.albumDir)
	if err != nil {
		return err
	}

	clear(f.cache)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			file, err := os.Open(filepath.Join(f.albumDir, entry.Name()))
			if err != nil {
				return err
			}
			defer file.Close()

			var album Album
			if err := json.NewDecoder(file).Decode(&album); err != nil {
				return err
			}
			f.cache[album.key] = album
		}
	}
	return nil
}

// // func (d *localFileSystem) sanatizeFileName(filename string) string {
// // 	sanitizer := strings.NewReplacer(
// // 		"<", "",
// // 		">", "",
// // 		":", "",
// // 		"\"", "",
// // 		"/", "",
// // 		"\\", "",
// // 		"|", "",
// // 		"?", "",
// // 		"*", "",
// // 		"~", "",
// // 	)
// // 	return sanitizer.Replace(filename)
// // }
