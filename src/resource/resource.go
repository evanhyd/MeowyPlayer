package resource

import "path/filepath"

func GetImageBasePath() string {
	return "images"
}

func GetMusicBasePath() string {
	return "music"
}

func GetAlbumBasePath() string {
	return "album"
}

func GetImagePath(image string) string {
	return filepath.Join(GetImageBasePath(), image)
}

func GetMusicPath(music string) string {
	return filepath.Join(GetMusicBasePath(), music)
}

func GetAlbumFolderPath(album string) string {
	return filepath.Join(GetAlbumBasePath(), album)
}

func GetAlbumIconPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(album), "icon.png")
}

func GetAlbumConfigPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(album), "config.txt")
}
