package model

import "io"

// Storage interface that controls the access to album and music.
// All the storage functions must be concurrency safe.
// But does not guarantee atomic transaction or execution order.
type Storage interface {
	getAllAlbums() ([]Album, error)
	getAlbum(AlbumKey) (Album, error)
	uploadAlbum(Album) error
	removeAlbum(AlbumKey) error
	getMusic(MusicKey) (io.ReadSeekCloser, error)
	uploadMusic(Music, io.Reader) error
	removeMusic(MusicKey) error
}
