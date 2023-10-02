package client

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/json"
	"meowyplayer.com/utility/pattern"
)

var collectionData pattern.Data[resource.Collection]
var albumData pattern.Data[resource.Album]
var playListData pattern.Data[resource.PlayList]

// the album pointer parameter may refer to a temporary object from the view list
// we need the original one from the collection
func getSourceAlbum(album *resource.Album) *resource.Album {
	index := slices.IndexFunc(collectionData.Get().Albums, func(a resource.Album) bool { return a.Title == album.Title })
	return &collectionData.Get().Albums[index]
}

func reloadCollection() error {
	if err := json.Write(resource.CollectionPath(), collectionData.Get()); err != nil {
		return err
	}
	collection, err := LoadFromLocalCollection()
	if err != nil {
		return err
	}
	collectionData.Set(&collection)
	return nil
}

func reloadAlbum() error {
	albumData.Set(getSourceAlbum(albumData.Get()))
	return nil
}

func LoadFromLocalCollection() (resource.Collection, error) {
	inUse := resource.Collection{}
	if err := json.Read(resource.CollectionPath(), &inUse); err != nil {
		return inUse, err
	}

	for i := range inUse.Albums {
		inUse.Albums[i].Cover = resource.GetCover(&inUse.Albums[i])
	}

	return inUse, nil
}

func GetCollectionData() *pattern.Data[resource.Collection] {
	return &collectionData
}

func GetAlbumData() *pattern.Data[resource.Album] {
	return &albumData
}

func GetPlayListData() *pattern.Data[resource.PlayList] {
	return &playListData
}

func AddAlbum() error {
	inUse := collectionData.Get()

	//generate title
	title := ""
	for i := 0; i < math.MaxInt; i++ {
		title = fmt.Sprintf("Album (%v)", i)
		if !slices.ContainsFunc(inUse.Albums, func(a resource.Album) bool { return a.Title == title }) {
			break
		}
	}

	//generate album
	album := resource.Album{Date: time.Now(), Title: title}
	inUse.Albums = append(inUse.Albums, album)

	//generate album cover
	iconColor := color.NRGBA{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())}
	iconImage := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	iconImage.SetNRGBA(0, 0, iconColor)
	imageData := bytes.Buffer{}
	if err := png.Encode(&imageData, iconImage); err != nil {
		return err
	}
	if err := os.WriteFile(resource.CoverPath(&album), imageData.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return reloadCollection()
}

func estimateMP3DataLength(data []byte) time.Duration {
	decoder, err := mp3.NewDecoder(bytes.NewReader(data))
	assert.NoErr(err)
	seconds := float64(decoder.Length()) / float64(resource.SAMPLING_RATE) / float64(resource.NUM_OF_CHANNELS) / float64(resource.AUDIO_BIT_DEPTH)
	return time.Duration(seconds * float64(time.Second))
}

func AddLocalMusic(musicInfo fyne.URIReadCloser) error {
	music := resource.Music{Date: time.Now(), Title: musicInfo.URI().Name()}

	//copy the music file to the music repo
	data, err := os.ReadFile(musicInfo.URI().Path())
	if err != nil {
		return err
	}
	if err = os.WriteFile(resource.MusicPath(&music), data, os.ModePerm); err != nil {
		return err
	}
	music.Length = estimateMP3DataLength(data)

	album := getSourceAlbum(albumData.Get())
	album.MusicList = append(album.MusicList, music)
	if err := reloadCollection(); err != nil {
		return err
	}
	return reloadAlbum()
}

func DeleteAlbum(album *resource.Album) error {
	collection := collectionData.Get()
	index := slices.IndexFunc(collection.Albums, func(a resource.Album) bool { return a.Title == album.Title })
	last := len(collection.Albums) - 1

	//remove album icon
	if err := os.Remove(resource.CoverPath(album)); err != nil && !os.IsNotExist(err) {
		return err
	}

	//pop from the collection
	collection.Albums[index] = collection.Albums[last]
	collection.Albums = collection.Albums[:last]
	return reloadCollection()
}

func DeleteMusic(music *resource.Music) error {
	album := getSourceAlbum(albumData.Get())
	index := slices.IndexFunc(album.MusicList, func(m resource.Music) bool { return m.SimpleTitle() == music.SimpleTitle() })
	last := len(album.MusicList) - 1

	//pop form the album
	album.MusicList[index] = album.MusicList[last]
	album.MusicList = album.MusicList[:last]

	if err := reloadCollection(); err != nil {
		return err
	}
	return reloadAlbum()
}

func UpdateAlbumTitle(album *resource.Album, title string) error {
	if slices.ContainsFunc(collectionData.Get().Albums, func(a resource.Album) bool { return a.Title == title }) {
		return fmt.Errorf("album \"%v\" already exists", title)
	}

	//update timestamp
	collectionData.Get().Date = time.Now()
	source := getSourceAlbum(album)
	source.Date = time.Now()

	//rename the album cover
	oldPath := resource.CoverPath(source)
	source.Title = title
	if err := os.Rename(oldPath, resource.CoverPath(source)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return reloadCollection()
}

func UpdateAlbumCover(album *resource.Album, iconPath string) error {
	album = getSourceAlbum(album)

	//update timestamp
	album.Date = time.Now()
	collectionData.Get().Date = time.Now()

	//update cover image
	icon, err := os.ReadFile(iconPath)
	if err != nil {
		return err
	}
	if err = os.WriteFile(resource.CoverPath(album), icon, os.ModePerm); err != nil {
		return err
	}
	return reloadCollection()
}
