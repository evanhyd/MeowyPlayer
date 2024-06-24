package model

import (
	"playground/pattern"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

type Client struct {
	storage            FileSystem
	onAlbumsChanged    pattern.Subject[[]Album]
	onAlbumSelected    pattern.Subject[Album]
	onAlbumViewFocused pattern.Subject[bool]
	onMusicViewFocused pattern.Subject[bool]
}

var client Client

func CreateClient(fileSystem FileSystem) {
	client = Client{
		storage:            fileSystem,
		onAlbumsChanged:    pattern.MakeSubject[[]Album](),
		onAlbumSelected:    pattern.MakeSubject[Album](),
		onAlbumViewFocused: pattern.MakeSubject[bool](),
		onMusicViewFocused: pattern.MakeSubject[bool](),
	}
}

func GetClient() *Client {
	return &client
}

func (m *Client) Run() error {
	if err := m.storage.initialize(); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) GetAlbum(key AlbumKey) (Album, error) {
	return m.storage.getAlbum(key)
}

func (m *Client) CreateAlbum(title string, cover fyne.Resource) error {
	album := Album{key: AlbumKey(uuid.NewString()), date: time.Now()}
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.title = title
	album.cover = cover
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) RemoveAlbum(key AlbumKey) error {
	if err := m.storage.removeAlbum(key); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) RemoveMusicFromAlbum(key AlbumKey, mKey MusicKey) error {
	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.music = slices.DeleteFunc(album.music, func(m Music) bool { return m.Key() == mKey })
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	return m.notifyAlbumsChanges()
}

func (m *Client) SelectAlbum(key AlbumKey) error {
	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	m.onAlbumSelected.NotifyAll(album)
	m.onMusicViewFocused.NotifyAll(true)
	return nil
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

func (m *Client) notifyAlbumsChanges() error {
	albums, err := m.storage.getAllAlbums()
	if err != nil {
		return err
	}
	m.onAlbumsChanged.NotifyAll(albums)
	return nil
}
