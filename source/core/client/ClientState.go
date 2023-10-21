package client

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/json"
	"meowyplayer.com/utility/pattern"
)

var state *ClientState = NewClientState()

func GetInstance() *ClientState {
	return state
}

type ClientState struct {
	accessLock sync.Mutex
	collection pattern.Data[resource.Collection]
	albumEvent pattern.SubjectBase[resource.Album]
	albumKey   string //keep the key instead of a copy
	playList   pattern.Data[resource.PlayList]
}

func NewClientState() *ClientState {
	return &ClientState{}
}

func (s *ClientState) save() error {
	return json.WriteFile(resource.CollectionPath(), s.collection.Get())
}

func (s *ClientState) GetAlbum() resource.Album {
	return s.collection.Get().Albums[s.albumKey]
}

func (s *ClientState) SetCollection(collection resource.Collection) {
	s.collection.Set(collection)
}

func (s *ClientState) SetAlbum(album resource.Album) {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()
	if album, ok := s.collection.Get().Albums[album.Title]; ok {
		s.albumKey = album.Title
		s.albumEvent.NotifyAll(album)
	} else {
		log.Printf("setting invalid album in the client state: %v", album.Title)
	}
}

func (s *ClientState) SetPlayList(playList *resource.PlayList) {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()
	s.playList.Set(*playList)
}

func (s *ClientState) AddCollectionListener(observer pattern.Observer[resource.Collection]) {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()
	s.collection.Attach(observer)
}

func (s *ClientState) AddAlbumListener(observer pattern.Observer[resource.Album]) {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()
	s.albumEvent.Attach(observer)
}

func (s *ClientState) AddPlayListListener(observer pattern.Observer[resource.PlayList]) {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()
	s.playList.Attach(observer)
}

func (s *ClientState) AddAlbum(album resource.Album) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	//add the album to the collection, then refresh
	if _, exist := s.collection.Get().Albums[album.Title]; !exist {

		//add album icon to the icon path
		if err := os.WriteFile(resource.CoverPath(&album), album.Cover.Content(), 0777); err != nil {
			return err
		}

		//add album to the collection
		collection := s.collection.Get()
		collection.Date = time.Now()
		collection.Albums[album.Title] = album
		s.collection.Set(collection)
		return s.save()
	} else {
		return fmt.Errorf("failed to add the duplicated album: %v", album.Title)
	}
}

func (s *ClientState) DeleteAlbum(album resource.Album) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	log.Printf("delete %v\n", album.Title)

	//delete album icon
	if err := os.Remove(resource.CoverPath(&album)); err != nil && !os.IsNotExist(err) {
		return err
	}

	//delete album from the collection
	collection := s.collection.Get()
	collection.Date = time.Now()
	delete(collection.Albums, album.Title)
	s.collection.Set(collection)
	return s.save()
}

func (s *ClientState) UpdateAlbumTitle(targetAlbum resource.Album, newTitle string) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	log.Printf("rename %v to %v\n", targetAlbum.Title, newTitle)

	//target album must exist
	album, exist := s.collection.Get().Albums[targetAlbum.Title]
	if !exist {
		return fmt.Errorf("failed to rename the title of a non-existed album: %v", targetAlbum.Title)
	}

	//album with the new title must not already exist
	if _, exist := s.collection.Get().Albums[newTitle]; exist {
		return fmt.Errorf("album title already exists: %v", newTitle)
	}

	//add the new album to the collection
	collection := s.collection.Get()
	collection.Date = time.Now()
	album.Date = time.Now()
	album.Title = newTitle
	delete(collection.Albums, targetAlbum.Title)
	collection.Albums[newTitle] = album
	if s.albumKey == targetAlbum.Title {
		s.albumKey = newTitle
	}

	//rename the album cover
	if err := os.Rename(resource.CoverPath(&targetAlbum), resource.CoverPath(&album)); err != nil && !os.IsNotExist(err) {
		return err
	}

	//update collection date
	s.collection.Set(collection)
	return s.save()
}

func (s *ClientState) UpdateAlbumCover(album resource.Album, iconPath string) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	log.Printf("update %v's cover: %v\n", album.Title, iconPath)

	//target album must exist
	if _, exist := s.collection.Get().Albums[album.Title]; !exist {
		return fmt.Errorf("failed to update the cover of a non-existed album: %v", album.Title)
	}

	//update cover image
	icon, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		return err
	}
	if err = os.WriteFile(resource.CoverPath(&album), icon.Content(), os.ModePerm); err != nil {
		return err
	}

	//update timestamp
	collection := s.collection.Get()
	collection.Date = time.Now()
	album.Date = time.Now()
	album.Cover = icon
	collection.Albums[album.Title] = album
	s.collection.Set(collection)
	return s.save()
}

func (s *ClientState) AddMusic(music resource.Music, musicData []byte) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	//album must exist
	album, exist := s.collection.Get().Albums[s.albumKey]
	if !exist {
		return fmt.Errorf("failed to add the music to a non-existed album: %v", s.albumKey)
	}

	//write data to the music repo
	if err := os.WriteFile(resource.MusicPath(&music), musicData, 0777); err != nil {
		return err
	}

	//updaite collection date, album date, album music list
	collection := s.collection.Get()
	collection.Date = time.Now()
	album.Date = time.Now()
	album.MusicList[music.Title] = music
	collection.Albums[album.Title] = album
	s.collection.Set(collection)
	s.albumEvent.NotifyAll(album)
	return s.save()
}

func (s *ClientState) DeleteMusic(music resource.Music) error {
	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	//album must exist
	album, exist := s.collection.Get().Albums[s.albumKey]
	if !exist {
		return fmt.Errorf("failed to delete the music to a non-existed album: %v", s.albumKey)
	}

	//remove from the collection
	//but not delete it from the music repo
	collection := s.collection.Get()
	collection.Date = time.Now()
	album.Date = time.Now()
	delete(album.MusicList, music.Title)
	collection.Albums[album.Title] = album
	s.collection.Set(collection)
	s.albumEvent.NotifyAll(album)
	return s.save()
}
