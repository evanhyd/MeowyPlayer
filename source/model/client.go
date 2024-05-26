package model

import (
	"fmt"
	"playground/pattern"
	"slices"

	"fyne.io/fyne/v2"
)

type Client struct {
	fileSystem         FileSystem
	onAlbumsChanged    pattern.Subject[[]Album]
	onAlbumSelected    pattern.Subject[Album]
	onAlbumViewFocused pattern.Subject[bool]
	onMusicViewFocused pattern.Subject[bool]
}

func NewClient(fileSystem FileSystem) Client {
	return Client{
		fileSystem:         fileSystem,
		onAlbumsChanged:    pattern.MakeSubject[[]Album](),
		onAlbumSelected:    pattern.MakeSubject[Album](),
		onAlbumViewFocused: pattern.MakeSubject[bool](),
		onMusicViewFocused: pattern.MakeSubject[bool](),
	}
}

func (m *Client) Initialize() error {
	if err := m.fileSystem.initialize(); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) GetAlbum(key AlbumKey) Album {
	album, err := m.fileSystem.getAlbum(key)
	if err != nil {
		fyne.LogError(fmt.Sprintf("failed to get album by key %v", key), err)
	}
	return album
}

func (m *Client) CreateAlbum(title string, cover fyne.Resource) error {
	_, err := m.fileSystem.createAlbum(Album{title: title, cover: cover.Content()})
	if err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	album := m.GetAlbum(key)
	album.title = title
	album.cover = cover.Content()

	if err := m.fileSystem.updateAlbum(album); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) RemoveAlbum(key AlbumKey) error {
	if err := m.fileSystem.removeAlbum(key); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) RemoveMusicFromAlbum(aKey AlbumKey, mKey MusicKey) error {
	album := m.GetAlbum(aKey)
	album.music = slices.DeleteFunc(album.music, func(m Music) bool { return m.Key() == mKey })
	if err := m.fileSystem.updateAlbum(album); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) notifyAlbumsChanges() error {
	albums, err := m.fileSystem.getAllAlbums()
	if err != nil {
		return err
	}
	m.onAlbumsChanged.NotifyAll(albums)
	return nil
}

func (m *Client) SelectAlbum(album Album) {
	m.onAlbumSelected.NotifyAll(album)
	m.onMusicViewFocused.NotifyAll(true)
}

func (m *Client) FocusAlbumView() {
	m.onAlbumViewFocused.NotifyAll(true)
}

func (m *Client) OnAlbumsChanged() pattern.Subject[[]Album] {
	return m.onAlbumsChanged
}

func (m *Client) OnAlbumSelected() pattern.Subject[Album] {
	return m.onAlbumSelected
}

func (m *Client) OnAlbumViewFocused() pattern.Subject[bool] {
	return m.onAlbumViewFocused
}

func (m *Client) OnMusicViewFocused() pattern.Subject[bool] {
	return m.onMusicViewFocused
}
