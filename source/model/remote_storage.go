package model

import (
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
