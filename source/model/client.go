package model

import (
	"io"
	"playground/browser"
	"playground/pattern"
	"slices"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

type Client struct {
	sync.Mutex         //read -> modify -> upload back may have intervene sequence
	storage            Storage
	onStorageLoad      pattern.Subject[[]Album]
	onAlbumSelected    pattern.Subject[Album]
	onAlbumViewFocused pattern.Subject[bool]
	onMusicViewFocused pattern.Subject[bool]
}

var client Client

func InitClient(storage Storage) {
	client = Client{
		storage:            storage,
		onStorageLoad:      pattern.MakeSubject[[]Album](),
		onAlbumSelected:    pattern.MakeSubject[Album](),
		onAlbumViewFocused: pattern.MakeSubject[bool](),
		onMusicViewFocused: pattern.MakeSubject[bool](),
	}
}

func Instance() *Client {
	return &client
}

func (m *Client) OnStorageLoaded() pattern.Subject[[]Album] {
	return m.onStorageLoad
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

func (m *Client) Run() error {
	if err := m.storage.initialize(); err != nil {
		return err
	}
	return m.reloadStorage()
}

func (m *Client) GetAlbum(key AlbumKey) (Album, error) {
	m.Lock()
	defer m.Unlock()

	return m.storage.getAlbum(key)
}

func (m *Client) GetAllAlbums() ([]Album, error) {
	m.Lock()
	defer m.Unlock()

	return m.storage.getAllAlbums()
}

func (m *Client) CreateAlbum(title string, cover fyne.Resource) error {
	m.Lock()
	defer m.Unlock()

	album := Album{key: AlbumKey(uuid.NewString()), title: title, cover: cover}
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	return m.reloadStorage()
}

func (m *Client) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	m.Lock()
	defer m.Unlock()

	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.title = title
	album.cover = cover
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	return m.reloadStorage()
}

func (m *Client) RemoveAlbum(key AlbumKey) error {
	m.Lock()
	defer m.Unlock()

	if err := m.storage.removeAlbum(key); err != nil {
		return err
	}
	return m.reloadStorage()
}

func (m *Client) AddMusicToAlbum(key AlbumKey, result browser.Result, reader io.Reader) error {
	m.Lock()
	defer m.Unlock()

	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	music := Music{title: result.Title, length: result.Length, platform: result.Platform, id: result.VideoID}
	album.music = append(album.music, music)
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	m.storage.uploadMusic(music, reader)
	return m.reloadStorage()
}

func (m *Client) RemoveMusicFromAlbum(key AlbumKey, mKey MusicKey) error {
	m.Lock()
	defer m.Unlock()

	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.music = slices.DeleteFunc(album.music, func(m Music) bool { return m.Key() == mKey })
	if err := m.storage.uploadAlbum(album); err != nil {
		return err
	}
	return m.reloadStorage()
}

func (m *Client) GetMusic(key MusicKey) (io.ReadSeekCloser, error) {
	m.Lock()
	defer m.Unlock()
	return m.storage.getMusic(key)
}

func (m *Client) SelectAlbum(key AlbumKey) error {
	m.Lock()
	defer m.Unlock()

	album, err := m.storage.getAlbum(key)
	if err != nil {
		return err
	}
	m.onAlbumSelected.NotifyAll(album)
	m.onMusicViewFocused.NotifyAll(true)
	return nil
}

func (m *Client) reloadStorage() error {
	albums, err := m.storage.getAllAlbums()
	if err != nil {
		return err
	}
	m.onStorageLoad.NotifyAll(albums)
	return nil
}

func (m *Client) FocusAlbumView() {
	m.onAlbumViewFocused.NotifyAll(true)
}
