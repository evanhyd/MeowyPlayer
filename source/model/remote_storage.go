package model

import (
	"encoding/json"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

var _ Storage = &remoteStorage{}

type remoteStorage struct {
	localStorage
}

func newRemoteStorage() *remoteStorage {
	const kStorage = "storage"
	s := remoteStorage{
		localStorage{
			albumDir: filepath.Join(kStorage, "remote"),
			musicDir: filepath.Join(kStorage, "music"),
		},
	}
	if err := os.MkdirAll(s.albumDir, 0700); err != nil {
		fyne.LogError("can not create local storage album dir", err)
	}
	if err := os.MkdirAll(s.musicDir, 0700); err != nil {
		fyne.LogError("can not create music dir", err)
	}
	return &s
}

func (s *remoteStorage) getAllAlbums() ([]Album, error) {
	albums, err := NetworkClient().getAllAlbums()
	if err != nil {
		return nil, err
	}

	//remove old albums
	dirs, err := os.ReadDir(s.localStorage.albumDir)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if err := os.RemoveAll(filepath.Join(s.albumDir, dir.Name())); err != nil {
			return nil, err
		}
	}

	//get new albums
	for i := range albums {
		file, err := os.OpenFile(s.albumPath(albums[i].Key()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, nil
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(albums[i]); err != nil {
			return nil, nil
		}
	}

	return albums, nil
}

func (s *remoteStorage) uploadAlbum(album Album) error {
	if err := NetworkClient().uploadAlbum(album); err != nil {
		return err
	}
	return s.localStorage.uploadAlbum(album)
}

func (s *remoteStorage) removeAlbum(key AlbumKey) error {
	if err := NetworkClient().removeAlbum(key); err != nil {
		return err
	}
	return s.localStorage.removeAlbum(key)
}
