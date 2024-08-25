package model

import (
	"io"
	"meowyplayer/browser"
	"meowyplayer/util"
	"slices"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

type uiClient struct {
	sync.Mutex         //read -> modify -> upload back may have intervene sequence
	storage            Storage
	onStorageLoad      util.Subject[[]Album]
	onAlbumSelected    util.Subject[Album]
	onAlbumViewFocused util.Subject[bool]
	onMusicViewFocused util.Subject[bool]
}

var uiClientInstance uiClient

func UIClient() *uiClient {
	return &uiClientInstance
}

func InitUIClient() error {
	uiClientInstance = uiClient{
		onStorageLoad:      util.MakeSubject[[]Album](),
		onAlbumSelected:    util.MakeSubject[Album](),
		onAlbumViewFocused: util.MakeSubject[bool](),
		onMusicViewFocused: util.MakeSubject[bool](),
	}
	return nil
}

func (c *uiClient) reloadStorage() error {
	albums, err := c.storage.getAllAlbums()
	if err != nil {
		return err
	}
	c.onStorageLoad.NotifyAll(albums)
	return nil
}

func (c *uiClient) setStorage(storage Storage) error {
	c.Lock()
	defer c.Unlock()
	c.onAlbumSelected.NotifyAll(Album{})
	c.onAlbumViewFocused.NotifyAll(true)
	c.storage = storage
	return c.reloadStorage()
}

func (c *uiClient) GetAlbum(key AlbumKey) (Album, error) {
	c.Lock()
	defer c.Unlock()
	return c.storage.getAlbum(key)
}

func (c *uiClient) GetAllAlbums() ([]Album, error) {
	c.Lock()
	defer c.Unlock()
	return c.storage.getAllAlbums()
}

func (c *uiClient) CreateAlbum(title string, cover fyne.Resource) error {
	c.Lock()
	defer c.Unlock()
	album := Album{key: newRandomAlbumKey(), date: time.Now(), title: title, cover: cover}
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *uiClient) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	c.Lock()
	defer c.Unlock()
	album, err := c.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.date = time.Now()
	album.title = title
	album.cover = cover
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *uiClient) RemoveAlbum(key AlbumKey) error {
	c.Lock()
	defer c.Unlock()
	if err := c.storage.removeAlbum(key); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *uiClient) AddMusicToAlbum(key AlbumKey, result browser.Result, reader io.Reader) error {
	c.Lock()
	defer c.Unlock()
	album, err := c.storage.getAlbum(key)
	if err != nil {
		return err
	}
	timestamp := time.Now()
	music := Music{date: timestamp, title: result.Title, length: result.Length, platform: result.Platform, id: result.VideoID}
	album.date = timestamp
	album.music = append(album.music, music)
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	if err := c.storage.uploadMusic(music, reader); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *uiClient) RemoveMusicFromAlbum(key AlbumKey, mKey MusicKey) error {
	c.Lock()
	defer c.Unlock()

	album, err := c.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.date = time.Now()
	album.music = slices.DeleteFunc(album.music, func(m Music) bool { return m.Key() == mKey }) //consider map
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *uiClient) GetMusic(key MusicKey) (io.ReadSeekCloser, error) {
	c.Lock()
	defer c.Unlock()
	return c.storage.getMusic(key)
}

func (c *uiClient) SelectAlbum(key AlbumKey) error {
	c.Lock()
	defer c.Unlock()
	album, err := c.storage.getAlbum(key)
	if err != nil {
		return err
	}
	c.onAlbumSelected.NotifyAll(album)
	c.onMusicViewFocused.NotifyAll(true)
	return nil
}

func (c *uiClient) FocusAlbumView() {
	c.onAlbumViewFocused.NotifyAll(true)
}

func (c *uiClient) OnStorageLoaded() util.Subject[[]Album] {
	return c.onStorageLoad
}

func (c *uiClient) OnAlbumSelected() util.Subject[Album] {
	return c.onAlbumSelected
}

func (c *uiClient) OnAlbumViewFocused() util.Subject[bool] {
	return c.onAlbumViewFocused
}

func (c *uiClient) OnMusicViewFocused() util.Subject[bool] {
	return c.onMusicViewFocused
}
