package model

import (
	"io"
	"meowyplayer/scraper"
	"meowyplayer/util"
	"os"
	"slices"
	"time"

	"fyne.io/fyne/v2"
)

type ViewID int

const (
	KAlbumView ViewID = iota
	KMusicView
)

type storageClient struct {
	storage              Storage
	onStorageLoad        util.Subject[[]Album]
	onAlbumSelected      util.Subject[Album]
	onViewFocused        util.Subject[ViewID]
	onMusicSyncActivated util.Subject[bool]
}

var storageClientInstance storageClient

func StorageClient() *storageClient {
	return &storageClientInstance
}

func InitStorageClient() error {
	const kStorage = "storage"
	storageClientInstance = storageClient{
		onStorageLoad:        util.MakeSubject[[]Album](),
		onAlbumSelected:      util.MakeSubject[Album](),
		onViewFocused:        util.MakeSubject[ViewID](),
		onMusicSyncActivated: util.MakeSubject[bool](),
	}
	return os.MkdirAll(kStorage, 0700)
}

func (c *storageClient) reloadStorage() error {
	albums, err := c.GetAllAlbums()
	if err != nil {
		return err
	}
	c.onStorageLoad.NotifyAll(albums)
	return nil
}

func (c *storageClient) setStorage(storage Storage) error {
	c.storage = storage
	c.onAlbumSelected.NotifyAll(Album{})
	c.onViewFocused.NotifyAll(KAlbumView)
	return c.reloadStorage()
}

func (c *storageClient) GetAlbum(key AlbumKey) (Album, error) {
	return c.storage.getAlbum(key)
}

func (c *storageClient) GetAllAlbums() ([]Album, error) {
	return c.storage.getAllAlbums()
}

func (c *storageClient) UploadAlbum(album Album) error {
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *storageClient) CreateAlbum(title string, cover fyne.Resource) error {
	return c.UploadAlbum(Album{key: newRandomAlbumKey(), date: time.Now(), title: title, cover: cover})
}

func (c *storageClient) EditAlbum(key AlbumKey, title string, cover fyne.Resource) error {
	album, err := c.GetAlbum(key)
	if err != nil {
		return err
	}
	album.date = time.Now()
	album.title = title
	album.cover = cover
	return c.UploadAlbum(album)
}

func (c *storageClient) RemoveAlbum(key AlbumKey) error {
	if err := c.storage.removeAlbum(key); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *storageClient) SelectAlbum(key AlbumKey) error {
	album, err := c.GetAlbum(key)
	if err != nil {
		return err
	}
	c.onAlbumSelected.NotifyAll(album)
	c.onViewFocused.NotifyAll(KMusicView)
	return nil
}

func (c *storageClient) GetMusic(key MusicKey) (io.ReadSeekCloser, error) {
	return c.storage.getMusic(key)
}

func (c *storageClient) SyncMusic(result scraper.Result) error {
	c.onMusicSyncActivated.NotifyAll(true) // UI update callback
	defer c.onMusicSyncActivated.NotifyAll(false)

	content, err := scraper.NewYouTubeDownloader().Download(&result)
	if err != nil {
		return err
	}
	defer content.Close()
	return c.storage.uploadMusic(Music{platform: result.Platform, id: result.ID}, content)
}

func (c *storageClient) UploadMusicToAlbum(key AlbumKey, result scraper.Result) error {
	album, err := c.storage.getAlbum(key)
	if err != nil {
		return err
	}
	album.date = time.Now()
	album.music = append(album.music, Music{
		date:     album.date,
		title:    result.Title,
		length:   result.Length,
		platform: result.Platform,
		id:       result.ID},
	)
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *storageClient) RemoveMusicFromAlbum(albumKey AlbumKey, musicKey MusicKey) error {
	album, err := c.storage.getAlbum(albumKey)
	if err != nil {
		return err
	}
	album.date = time.Now()
	album.music = slices.DeleteFunc(album.music, func(m Music) bool { return m.Key() == musicKey }) //consider map
	if err := c.storage.uploadAlbum(album); err != nil {
		return err
	}
	return c.reloadStorage()
}

func (c *storageClient) FocusAlbumView() {
	c.onViewFocused.NotifyAll(KAlbumView)
}

func (c *storageClient) OnStorageLoaded() util.Subject[[]Album] {
	return c.onStorageLoad
}

func (c *storageClient) OnAlbumSelected() util.Subject[Album] {
	return c.onAlbumSelected
}

func (c *storageClient) OnViewFocused() util.Subject[ViewID] {
	return c.onViewFocused
}

func (c *storageClient) OnMusicSyncActivated() util.Subject[bool] {
	return c.onMusicSyncActivated
}
