package model

import (
	"fmt"
	"playground/pattern"
	"slices"

	"fyne.io/fyne/v2"
)

type MusicClient struct {
	fileSystem         FileSystem
	onAlbumsChanged    pattern.Subject[[]Album]
	onAlbumSelected    pattern.Subject[Album]
	onAlbumViewFocused pattern.Subject[bool]
	onMusicViewFocused pattern.Subject[bool]
}

func NewStorageClient(fileSystem FileSystem) MusicClient {
	return MusicClient{
		fileSystem:         fileSystem,
		onAlbumsChanged:    pattern.MakeSubject[[]Album](),
		onAlbumSelected:    pattern.MakeSubject[Album](),
		onAlbumViewFocused: pattern.MakeSubject[bool](),
		onMusicViewFocused: pattern.MakeSubject[bool](),
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
	album := Album{title: title, cover: cover.Content()}
	if _, err := m.fileSystem.createAlbum(album); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *MusicClient) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	album := m.GetAlbum(key)
	album.title = title
	album.cover = cover.Content()

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

func (m *MusicClient) RemoveMusic(aKey AlbumKey, mKey MusicKey) error {
	album := m.GetAlbum(aKey)
	album.music = slices.DeleteFunc(album.music, func(m Music) bool { return m.Key() == mKey })
	if err := m.fileSystem.updateAlbum(album); err != nil {
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

func (m *MusicClient) SelectAlbum(album Album) {
	m.onAlbumSelected.NotifyAll(album)
	m.onMusicViewFocused.NotifyAll(true)
}

func (m *MusicClient) FocusAlbumView() {
	m.onAlbumViewFocused.NotifyAll(true)
}

func (m *MusicClient) OnAlbumsChanged() pattern.Subject[[]Album] {
	return m.onAlbumsChanged
}

func (m *MusicClient) OnAlbumSelected() pattern.Subject[Album] {
	return m.onAlbumSelected
}

func (m *MusicClient) OnAlbumViewFocused() pattern.Subject[bool] {
	return m.onAlbumViewFocused
}

func (m *MusicClient) OnMusicViewFocused() pattern.Subject[bool] {
	return m.onMusicViewFocused
}
