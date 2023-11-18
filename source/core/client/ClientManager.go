package client

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	"meowyplayer.com/core/player"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/logger"
	"meowyplayer.com/utility/network/fileformat"
	"meowyplayer.com/utility/pattern"
	"meowyplayer.com/utility/ujson"
)

var manager = clientManager{}

func Manager() *clientManager {
	return &manager
}

type clientManager struct {
	accessLock sync.Mutex
	collection pattern.Data[resource.Collection]
	albumTitle string //key in the collection
	albumEvent pattern.SubjectBase[resource.Album]
	playList   pattern.Data[player.PlayList]
}

func (c *clientManager) Initialize() error {
	_, err := os.Stat(resource.CollectionFile())
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	//create default collection
	if os.IsNotExist(err) {
		collection := resource.Collection{Date: time.Now(), Albums: make(map[string]resource.Album)}
		if err := ujson.WriteFile(resource.CollectionFile(), collection); err != nil {
			return err
		}
	}
	return c.load()
}

func (c *clientManager) save() error {
	return ujson.WriteFile(resource.CollectionFile(), c.collection.Get())
}

func (c *clientManager) load() error {
	collection := resource.Collection{}
	if err := ujson.ReadFile(resource.CollectionFile(), &collection); err != nil {
		return err
	}

	for title, album := range collection.Albums {
		album.Cover = resource.Cover(&album)
		collection.Albums[title] = album
	}

	c.collection.Set(collection)
	return nil
}

func (c *clientManager) TempCollection() resource.Collection {
	return c.collection.Get()
}

func (c *clientManager) Album() resource.Album {
	return c.collection.Get().Albums[c.albumTitle]
}

func (c *clientManager) SetAlbum(album resource.Album) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	if album, ok := c.collection.Get().Albums[album.Title]; ok {
		c.albumTitle = album.Title
		c.albumEvent.NotifyAll(album)
	} else {
		logger.Error(fmt.Errorf("setting invalid album - %v", album.Title), 0)
	}
}

func (c *clientManager) SetPlayList(playList *player.PlayList) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	c.playList.Set(*playList)
}

func (c *clientManager) AddCollectionListener(observer pattern.Observer[resource.Collection]) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	c.collection.Attach(observer)
}

func (c *clientManager) AddAlbumListener(observer pattern.Observer[resource.Album]) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	c.albumEvent.Attach(observer)
}

func (c *clientManager) AddPlayListListener(observer pattern.Observer[player.PlayList]) {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()
	c.playList.Attach(observer)
}

func (c *clientManager) addAlbum(album resource.Album) error {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()

	//add the album to the collection, then refresh
	if _, exist := c.collection.Get().Albums[album.Title]; !exist {

		//add album icon to the icon path
		if err := os.WriteFile(resource.CoverPath(&album), album.Cover.Content(), 0777); err != nil {
			return err
		}

		//add album to the collection
		collection := c.collection.Get()
		collection.Date = time.Now()
		collection.Albums[album.Title] = album
		c.collection.Set(collection)
		return c.save()
	} else {
		return fmt.Errorf("failed to add the duplicated album: %v", album.Title)
	}
}

func (c *clientManager) AddRandomAlbum() error {
	//generate album cover
	iconColor := color.NRGBA{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())}
	iconImage := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	iconImage.SetNRGBA(0, 0, iconColor)
	imageData := bytes.Buffer{}
	if err := png.Encode(&imageData, iconImage); err != nil {
		return err
	}

	//generate album
	album := resource.Album{
		Date:      time.Now(),
		Title:     "",
		MusicList: make(map[string]resource.Music),
		Cover:     fyne.NewStaticResource("", imageData.Bytes()),
	}

	// try 100 possible titles until it fits in
	for i := 0; i < 100; i++ {
		album.Title = fmt.Sprintf("Album %v", i)
		if err := c.addAlbum(album); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to add new album")
}

func (c *clientManager) UpdateAlbumTitle(toRename resource.Album, newTitle string) error {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()

	newTitle = resource.SanatizeFileName(newTitle)
	toRename.Title = resource.SanatizeFileName(toRename.Title)
	log.Printf("rename %v to %v\n", toRename.Title, newTitle)

	//target album must exist
	album, exist := c.collection.Get().Albums[toRename.Title]
	if !exist {
		return fmt.Errorf("failed to rename the title of a non-existed album: %v", toRename.Title)
	}

	//album with the new title must not already exist
	if _, exist := c.collection.Get().Albums[newTitle]; exist {
		return fmt.Errorf("album title already exists: %v", newTitle)
	}

	//add the new album to the collection
	album.Date = time.Now()
	album.Title = newTitle
	collection := c.collection.Get()
	collection.Date = album.Date
	delete(collection.Albums, toRename.Title)
	collection.Albums[newTitle] = album
	c.collection.Set(collection)

	//update the album key to the new one
	//so the reference is not broken
	if c.albumTitle == toRename.Title {
		c.albumTitle = newTitle
	}

	//rename the album cover
	if err := os.Rename(resource.CoverPath(&toRename), resource.CoverPath(&album)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return c.save()
}

func (c *clientManager) UpdateAlbumCover(album resource.Album, iconPath string) error {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()

	log.Printf("update %v cover: %v\n", album.Title, iconPath)

	//target album must exist
	if _, exist := c.collection.Get().Albums[album.Title]; !exist {
		return fmt.Errorf("failed to update the cover of a non-existed album: %v", album.Title)
	}

	//update cover image
	icon, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		return err
	}
	if err = os.WriteFile(resource.CoverPath(&album), icon.Content(), 0777); err != nil {
		return err
	}

	//update timestamp
	album.Date = time.Now()
	album.Cover = icon
	collection := c.collection.Get()
	collection.Date = album.Date
	collection.Albums[album.Title] = album
	c.collection.Set(collection)
	return c.save()
}

func (c *clientManager) DeleteAlbum(album resource.Album) error {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()

	album.Title = resource.SanatizeFileName(c.albumTitle)
	log.Println("delete", album.Title)

	//delete album icon
	if err := os.RemoveAll(resource.CoverPath(&album)); err != nil {
		return err
	}

	//delete album from the collection
	collection := c.collection.Get()
	collection.Date = time.Now()
	delete(collection.Albums, album.Title)
	c.collection.Set(collection)
	return c.save()
}

func (c *clientManager) addMusic(music resource.Music, musicData []byte) error {
	c.accessLock.Lock()
	defer c.accessLock.Unlock()

	//album must exist
	album, exist := c.collection.Get().Albums[c.albumTitle]
	if !exist {
		return fmt.Errorf("failed to add the music to a non-existed album: %v", c.albumTitle)
	}

	//write data to the music repo
	if err := os.WriteFile(resource.MusicPath(&music), musicData, 0777); err != nil {
		return err
	}

	//updaite album date, album music list, collection date
	album.Date = time.Now()
	album.MusicList[music.Title] = music
	collection := c.collection.Get()
	collection.Date = album.Date
	collection.Albums[album.Title] = album
	c.collection.Set(collection)
	c.albumEvent.NotifyAll(album)
	return c.save()
}

func (c *clientManager) AddMusicFromDownloader(videoResult *fileformat.VideoResult, musicData []byte) error {
	music := resource.Music{
		Date:     time.Now(),
		Title:    resource.SanatizeFileName(videoResult.Title) + ".mp3",
		Length:   videoResult.Length,
		Platform: videoResult.Platform,
		ID:       videoResult.VideoID,
	}
	return c.addMusic(music, musicData)
}

func (c *clientManager) AddMusicFromURIReader(musicInfo fyne.URIReadCloser) error {
	estimateMP3DataLength := func(data []byte) (time.Duration, error) {
		decoder, err := mp3.NewDecoder(bytes.NewReader(data))
		if err != nil {
			return 0, err
		}
		seconds := float64(decoder.Length()) / float64(resource.SAMPLING_RATE) / float64(resource.NUM_OF_CHANNELS) / float64(resource.AUDIO_BIT_DEPTH)
		return time.Duration(seconds * float64(time.Second)), nil
	}

	musicData, err := os.ReadFile(musicInfo.URI().Path())
	if err != nil {
		return err
	}
	length, err := estimateMP3DataLength(musicData)
	if err != nil {
		return err
	}
	music := resource.Music{Date: time.Now(), Title: musicInfo.URI().Name(), Length: length}
	return c.addMusic(music, musicData)
}

func (s *clientManager) DeleteMusic(music resource.Music) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	//target album must exist
	album, exist := s.collection.Get().Albums[s.albumTitle]
	if !exist {
		return fmt.Errorf("failed to delete the music to a non-existed album: %v", s.albumTitle)
	}

	//remove from the collection
	//but not delete it from the music repo
	album.Date = time.Now()
	delete(album.MusicList, music.Title)
	collection := s.collection.Get()
	collection.Date = album.Date
	collection.Albums[album.Title] = album
	s.collection.Set(collection)
	s.albumEvent.NotifyAll(album)
	return s.save()
}