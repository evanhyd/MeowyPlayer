package client

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/pattern"
	"meowyplayer.com/utility/ujson"
)

var manager = clientManager{}

func Manager() *clientManager {
	return &manager
}

type clientManager struct {
	currentCollection     pattern.Data[resource.Collection]
	focusedAlbum          string
	onFocusedAlbumChanged pattern.SubjectBase[resource.Album] //callback when update the focused album
	onAlbumPlayed         pattern.SubjectBase[resource.Album] //callback when load the playlist from the album
}

/*
Initialize the client and load the collection config file.
Creates a default collection config file if not exists.
*/
func (c *clientManager) Initialize() error {
	if _, err := os.Stat(resource.CollectionFile()); errors.Is(err, fs.ErrNotExist) {
		collection := resource.Collection{Date: time.Now(), Albums: make(map[string]resource.Album)}
		if err := ujson.Write(resource.CollectionFile(), collection); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return c.load()
}

/*
Save the current collection to the collection config file.
*/
func (c *clientManager) save() error {
	return ujson.Write(resource.CollectionFile(), c.currentCollection.Get())
}

/*
Load from the collection config file to the current collection.
*/
func (c *clientManager) load() error {
	collection := resource.Collection{}
	if err := ujson.Read(resource.CollectionFile(), &collection); err != nil {
		return err
	}

	for title, album := range collection.Albums {
		album.Cover = resource.Cover(&album)
		collection.Albums[title] = album
	}

	c.currentCollection.Set(collection)
	return nil
}

/*
Return the current album.
*/
func (c *clientManager) Album() resource.Album {
	return c.currentCollection.Get().Albums[c.focusedAlbum]
}

/*
Set the current album.
*/
func (c *clientManager) SetAlbum(album resource.Album) error {
	if source, ok := c.currentCollection.Get().Albums[album.Title]; ok {
		c.focusedAlbum = source.Title
		c.onFocusedAlbumChanged.NotifyAll(source)
		return nil
	}
	return fmt.Errorf("setting invalid album - %v", album.Title)
}

/*
Add collection on-change event listener.
*/
func (c *clientManager) AddCollectionListener(observer pattern.Observer[resource.Collection]) {
	c.currentCollection.Attach(observer)
}

/*
Add on focused album changed listener.
*/
func (c *clientManager) AddAlbumListener(observer pattern.Observer[resource.Album]) {
	c.onFocusedAlbumChanged.Attach(observer)
}

/*
Add on album played listener.
*/
func (c *clientManager) AddAlbumPlayedListener(observer pattern.Observer[resource.Album]) {
	c.onAlbumPlayed.Attach(observer)
}

/*
Add album to the current collection.
Duplicated album title is prohibited.
Immediately save to the collection config file.
*/
func (c *clientManager) addAlbum(album resource.Album) error {
	collection := c.currentCollection.Get()
	album.Title = resource.SanatizeFileName(album.Title)
	album.Date = time.Now()

	//check for duplicted album title
	if _, exist := collection.Albums[album.Title]; !exist {

		//save icon
		if err := os.WriteFile(resource.CoverPath(&album), album.Cover.Content(), 0777); err != nil {
			return err
		}

		//add to the collection
		collection.Date = album.Date
		collection.Albums[album.Title] = album
		c.currentCollection.Set(collection)
		return c.save()
	}

	return fmt.Errorf("failed to add the duplicated album: %v", album.Title)
}

/*
Delete the album from the collection.
Immediately save to the collection config file.
*/
func (c *clientManager) DeleteAlbum(album resource.Album) error {
	album.Title = resource.SanatizeFileName(album.Title)

	//delete the icon
	if err := os.RemoveAll(resource.CoverPath(&album)); err != nil {
		return err
	}

	//delete from the collection
	collection := c.currentCollection.Get()
	delete(collection.Albums, album.Title)
	collection.Date = time.Now()
	c.currentCollection.Set(collection)
	return c.save()
}

/*
Update album's title to title.
The title must not already exists in the collection.
Immediately save to the collection config file.
*/
func (c *clientManager) UpdateAlbumTitle(album resource.Album, title string) error {
	collection := c.currentCollection.Get()

	title = resource.SanatizeFileName(title)
	if _, exist := collection.Albums[title]; exist {
		return fmt.Errorf("attempted to rename to an existed title: %v", title)
	}

	renamed, exist := collection.Albums[album.Title]
	if !exist {
		return fmt.Errorf("attempted to rename an invalid album: %v", album.Title)
	}
	renamed.Title = title
	renamed.Date = time.Now()

	//rename the cover
	if err := os.Rename(resource.CoverPath(&album), resource.CoverPath(&renamed)); err != nil {
		return err
	}
	delete(collection.Albums, album.Title)

	//update the collection
	collection.Date = renamed.Date
	collection.Albums[renamed.Title] = renamed
	c.currentCollection.Set(collection)

	//update the current album if necessary
	if c.focusedAlbum == album.Title {
		c.focusedAlbum = renamed.Title
		//TODO: notify listeners
	}

	return c.save()
}

/*
Update album's cover to the icon specified by the iconPath.
Immediately save to the collection config file.
*/
func (c *clientManager) UpdateAlbumCover(album resource.Album, iconPath string) error {
	collection := c.currentCollection.Get()

	source, exist := collection.Albums[album.Title]
	if !exist {
		return fmt.Errorf("attempted to update an invalid album's icon: %v", album.Title)
	}

	//update cover image
	icon, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		return err
	}
	if err = os.WriteFile(resource.CoverPath(&source), icon.Content(), 0777); err != nil {
		return err
	}
	source.Date = time.Now()
	source.Cover = icon

	//update the collection
	collection.Date = source.Date
	collection.Albums[source.Title] = source
	c.currentCollection.Set(collection)
	return c.save()
}

/*
Add music to the album.
Save musicData to the music directory.
Immediately save the collection config file.
*/
func (c *clientManager) addMusic(album resource.Album, music resource.Music, musicData []byte) error {
	collection := c.currentCollection.Get()

	source, exist := collection.Albums[album.Title]
	if !exist {
		return fmt.Errorf("attempted to add music to an invalid album: %v", album.Title)
	}
	source.Date = time.Now()
	source.MusicList[music.Title] = music

	//save the music data
	if err := os.WriteFile(resource.MusicPath(&music), musicData, 0777); err != nil {
		return err
	}

	//save the collection
	collection.Date = source.Date
	collection.Albums[source.Title] = source
	c.currentCollection.Set(collection)
	return c.save()
}

/*
Delete the music from the album.
Immediately save the collection config file.
*/
func (s *clientManager) DeleteMusic(album resource.Album, music resource.Music) error {
	collection := s.currentCollection.Get()

	source, exist := collection.Albums[album.Title]
	if !exist {
		return fmt.Errorf("attempted to delete a music from an invalid album: %v", album.Title)
	}
	source.Date = time.Now()
	delete(source.MusicList, music.Title)

	//remove from the collection
	//but not delete it from the music repo
	collection.Date = source.Date
	collection.Albums[source.Title] = source
	s.currentCollection.Set(collection)
	return s.save()
}
