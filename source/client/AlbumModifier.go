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

	"meowyplayer.com/source/resource"
)

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

	return reloadCollectionData()
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
	return reloadCollectionData()
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
	return reloadCollectionData()
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
	return reloadCollectionData()
}
