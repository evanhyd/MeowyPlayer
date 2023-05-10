package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

func AddMusicToRepository(videoID string, album player.Album, musicTitle string) error {
	const (
		y2mateServerUrl    = `https://www.y2mate.com/mates/analyze/ajax`
		y2mateConverterUrl = `https://www.y2mate.com/mates/convert`
		youtubeUrl         = `https://www.youtube.com/watch?v=`
	)

	//sanitize music title
	sanitizer := strings.NewReplacer(
		"<", "",
		">", "",
		":", "",
		"\"", "",
		"/", "",
		"\\", "",
		"|", "",
		"?", "",
		"*", "",
	)
	musicTitle = sanitizer.Replace(musicTitle) + ".mp3"

	//server POST
	youtubeVideoUrl := youtubeUrl + videoID
	serverQuery := url.Values{"url": {youtubeVideoUrl}, "q_auto": {"1"}, "ajax": {"1"}}.Encode()
	log.Printf("Video url: %v\n", youtubeVideoUrl)
	log.Printf("POST url: %v?%v\n", y2mateServerUrl, serverQuery)
	serverResp, err := http.Post(y2mateServerUrl, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(serverQuery))
	if err != nil {
		return err
	}
	defer serverResp.Body.Close()

	//scrape user ID
	serverBody, err := io.ReadAll(serverResp.Body)
	if err != nil {
		return err
	}
	userID := ""
	if begin := bytes.Index(serverBody, []byte(`var k__id = \"`)); begin != -1 {
		if end := bytes.Index(serverBody, []byte(`\"; var video_service =`)); end != -1 {
			userID = string(serverBody[begin+len(`var k__id = \"`) : end])
		}
	}
	if userID == "" {
		return errors.New("couldn't obtain user ID from server response")
	}
	log.Printf("User ID: %v\n", userID)

	//converter POST
	converterQuery := url.Values{"type": {"youtube"}, "_id": {userID}, "v_id": {videoID}, "ajax": {"1"}, "token": {""}, "ftype": {"mp3"}, "fquality": {"128"}}.Encode()
	converterResp, err := http.Post(y2mateConverterUrl, `application/x-www-form-urlencoded; charset=UTF-8`, strings.NewReader(converterQuery))
	if err != nil {
		return err
	}
	defer converterResp.Body.Close()

	//scrape file url
	converterBody, err := io.ReadAll(converterResp.Body)
	if err != nil {
		return err
	}
	fileUrl := ""
	if begin := bytes.Index(converterBody, []byte(`<a href=\"`)); begin != -1 {
		if end := bytes.Index(converterBody, []byte(`\" rel=\"nofollow\"`)); end != -1 {
			fileUrl = strings.ReplaceAll(string(converterBody[begin+len(`<a href=\"`):end]), `\/`, `/`)
		}
	}
	if fileUrl == "" {
		return errors.New("couldn't obtain file url")
	}

	//download mp3 to music folder
	fileResp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer fileResp.Body.Close()

	musicPath := resource.GetMusicPath(musicTitle)
	file, err := os.Create(musicPath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytesRead, err := io.Copy(file, fileResp.Body)
	if err != nil {
		return err
	}
	if bytesRead == 0 {
		return fmt.Errorf("failed to download due to server issue")
	}
	log.Printf("File size: %.2v mb\n", float64(bytesRead)/1024.0/1024.0)

	return player.AddMusicToAlbum(album, musicPath, musicTitle)
}
