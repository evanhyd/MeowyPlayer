package panel

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/resource"
	"meowyplayer.com/src/seeker"
)

const (
	ALBUM_CARD_WIDTH         = 128.0
	ALBUM_CARD_HEIGHT        = 128.0
	UPLOAD_DIALOG_WIN_WIDTH  = 512.0
	UPLOAD_DIALOG_WIN_HEIGHT = 512.0
)

//check for substring ignore cases
func SatisfyAlbumInfo(query string, data *custom_canvas.AlbumInfo) bool {
	return strings.Contains(strings.ToLower(data.Title), strings.ToLower(query))
}

//Load albums from the directories to the album search list
func LoadAlbumFromDir(panelInfo *custom_canvas.PanelInfo) error {

	//read directories
	albumDirs, err := os.ReadDir(resource.GetAlbumBasePath())
	if err != nil {
		return err
	}

	//read to album list
	panelInfo.AlbumSearchList.ClearData()
	for _, albumDir := range albumDirs {
		if albumDir.IsDir() {
			//read music config
			configFile, err := os.Open(resource.GetAlbumConfigPath(albumDir.Name()))
			if err != nil {
				return err
			}
			defer configFile.Close()

			buf := make([]byte, 32*1024)
			sz, err := configFile.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			cnt := bytes.Count(buf[:sz], []byte{'\n'})
			panelInfo.AlbumSearchList.AddData(custom_canvas.AlbumInfo{Title: albumDir.Name(), MusicCount: int32(cnt)})
		}
	}

	panelInfo.AlbumSearchList.ResetSearch()
	return nil
}

//Load music from the the album config
func LoadMusicFromAlbum(panelInfo *custom_canvas.PanelInfo) error {

	//load music config
	configFile, err := os.Open(resource.GetAlbumConfigPath(panelInfo.SelectedAlbumInfo.Title))
	if err != nil {
		return err
	}
	defer configFile.Close()

	//clear the old music
	panelInfo.MusicSearchList.ClearData()

	//read music name from config file
	scanner := bufio.NewScanner(configFile)
	for line := 0; scanner.Scan(); line++ {

		//open music file
		musicPath := resource.GetMusicPath(scanner.Text())
		musicInfo, err := os.Stat(musicPath)
		if err != nil {
			return err
		}

		//add new music
		sec := int64(float64(musicInfo.Size()) * 11.024576435167347 / float64(seeker.SAMPLING_RATE*seeker.AUDIO_BIT_DEPTH*seeker.NUM_OF_CHANNELS))
		panelInfo.MusicSearchList.AddData(custom_canvas.MusicInfo{Title: scanner.Text(), Duration: sec, Index: 0})
	}

	//sort and index the music
	sort.Slice(panelInfo.MusicSearchList.DataList, func(i, j int) bool {
		return strings.ToLower(panelInfo.MusicSearchList.DataList[i].Title) < strings.ToLower(panelInfo.MusicSearchList.DataList[j].Title)
	})
	for i := range panelInfo.MusicSearchList.DataList {
		panelInfo.MusicSearchList.DataList[i].Index = i
	}
	panelInfo.MusicSearchList.ResetSearch()
	return nil
}

//Create an album folder given by the image and title
func AddAlbum(imagePath, albumTitle string) error {

	albumPath := resource.GetAlbumFolderPath(albumTitle)

	if _, err := os.Stat(albumPath); os.IsNotExist(err) {

		//create album folder
		err := os.Mkdir(albumPath, os.ModePerm)
		if err != nil {
			return err
		}

		//copy album icon
		imageByte, err := os.ReadFile(imagePath)
		if err != nil {
			return err
		}
		err = os.WriteFile(resource.GetAlbumIconPath(albumTitle), imageByte, os.ModePerm)
		if err != nil {
			return err
		}

		//create empty config
		_, err = os.Create(resource.GetAlbumConfigPath(albumTitle))
		if err != nil {
			return err
		}
		return nil

	} else {
		return errors.New("duplicated album")
	}
}

func RemoveAlbum(albumTitle string) error {
	albumPath := resource.GetAlbumFolderPath(albumTitle)

	//remove album directory
	return os.RemoveAll(albumPath)
}
