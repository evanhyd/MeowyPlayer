package model

import (
	"fmt"
	"playground/pattern"
	"time"

	"fyne.io/fyne/v2"
)

type MusicClient struct {
	fileSystem      FileSystem
	focusedAlbum    Album
	onAlbumsChanged pattern.Subject[[]Album]
	onAlbumFocused  pattern.Subject[Album]
}

func NewStorageClient(fileSystem FileSystem) MusicClient {
	return MusicClient{
		fileSystem:      fileSystem,
		onAlbumsChanged: pattern.MakeSubject[[]Album](),
		onAlbumFocused:  pattern.MakeSubject[Album](),
	}
}

func (m *MusicClient) Initialize() error {
	if err := m.fileSystem.initialize(); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *MusicClient) GetAlbum(key AlbumKey) Album {
	album, err := m.fileSystem.getAlbum(key)
	if err != nil {
		fyne.LogError(fmt.Sprintf("failed to get album by key %v", key), err)
	}
	return album
}

func (m *MusicClient) CreateAlbum(title string, cover fyne.Resource) error {
	album := Album{date: time.Now(), title: title, cover: cover.Content()}
	if _, err := m.fileSystem.uploadAlbum(album); err != nil {
		return err
	}

	return m.notifyAlbumsChanges()
}

func (m *MusicClient) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	album := m.GetAlbum(key)
	album.title = title
	album.cover = cover.Content()
	album.date = time.Now()

	if err := m.fileSystem.updateAlbum(album); err != nil {
		return err
	}

	return m.notifyAlbumsChanges()
}

func (m *MusicClient) RemoveAlbum(key AlbumKey) error {
	if err := m.fileSystem.removeAlbum(key); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *MusicClient) notifyAlbumsChanges() error {
	albums, err := m.fileSystem.getAllAlbums()
	if err != nil {
		return err
	}
	m.onAlbumsChanged.NotifyAll(albums)
	return nil
}

func (m *MusicClient) RefreshAlbums() {
	m.notifyAlbumsChanges()
}

func (m *MusicClient) FocusAlbum(album Album) {
	m.focusedAlbum = album
	m.onAlbumFocused.NotifyAll(album)
}

func (m *MusicClient) OnAlbumsChanged() pattern.Subject[[]Album] {
	return m.onAlbumsChanged
}

func (m *MusicClient) OnAlbumFocused() pattern.Subject[Album] {
	return m.onAlbumFocused
}
