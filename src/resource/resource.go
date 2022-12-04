package resource

import (
	"path/filepath"
)

func GetBasePath() string {

	return ""
	// path, err := os.Executable()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// if strings.Contains(path, filepath.Join("var", "folders")) || strings.Contains(path, filepath.Join("Local", "Temp")) {
	// 	return ""
	// }

	// //full path is needed for native app
	// return filepath.Dir(path)
}

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
	return filepath.Join(GetBasePath(), GetImageBasePath(), image)
}

func GetMusicPath(music string) string {
	return filepath.Join(GetBasePath(), GetMusicBasePath(), music)
}

func GetAlbumFolderPath(album string) string {
	return filepath.Join(GetBasePath(), GetAlbumBasePath(), album)
}

func GetAlbumIconPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(album), "icon.png")
}

func GetAlbumConfigPath(album string) string {
	return filepath.Join(GetAlbumFolderPath(album), "config.txt")
}
