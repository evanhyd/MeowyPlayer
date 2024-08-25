package model

import "io"

type Storage interface {
	getAllAlbums() ([]Album, error)
	getAlbum(AlbumKey) (Album, error)
	getMusic(MusicKey) (io.ReadSeekCloser, error)
	uploadAlbum(Album) error
	uploadMusic(Music, io.Reader) error
	removeAlbum(AlbumKey) error
}
